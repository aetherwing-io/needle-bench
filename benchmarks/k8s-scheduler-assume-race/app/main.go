package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Simplified reproduction of k8s scheduler assume-and-bind race.
// The bug: assume() succeeds, but bind() fails asynchronously.
// handleError() re-queues the pod without checking if it was already bound
// by another goroutine, causing duplicate scheduling.

type Pod struct {
	Name     string
	NodeName string
	Bound    bool
}

type AssumeCache struct {
	mu       sync.Mutex
	assumed  map[string]*Pod
	bindings int64
}

func NewAssumeCache() *AssumeCache {
	return &AssumeCache{assumed: make(map[string]*Pod)}
}

// assume optimistically places a pod on a node
func (c *AssumeCache) assume(pod *Pod, node string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.assumed[pod.Name]; exists {
		return fmt.Errorf("pod %s already assumed", pod.Name)
	}
	pod.NodeName = node
	c.assumed[pod.Name] = pod
	return nil
}

// bind actually binds the pod (async, can fail)
func (c *AssumeCache) bind(pod *Pod) error {
	// Simulate occasional bind failure
	time.Sleep(time.Millisecond * 5)

	c.mu.Lock()
	defer c.mu.Unlock()

	assumed, exists := c.assumed[pod.Name]
	if !exists {
		return fmt.Errorf("pod %s not in assume cache", pod.Name)
	}

	// BUG: No atomic check between "is it bound?" and "mark it bound"
	// Another goroutine calling handleError can see Bound=false and re-queue
	assumed.Bound = true
	atomic.AddInt64(&c.bindings, 1)
	return nil
}

// handleError is called when scheduling fails — should re-queue the pod
// BUG: doesn't check if pod was already bound by async bind goroutine
func (c *AssumeCache) handleError(pod *Pod) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	assumed, exists := c.assumed[pod.Name]
	if !exists {
		return false
	}

	// BUG: This read of Bound races with bind() setting it to true.
	// In the real k8s code, this check exists but the window between
	// checking and re-queuing allows duplicate scheduling.
	if assumed.Bound {
		return false // already bound, don't re-queue
	}

	// Re-queue for scheduling — but bind() might complete right after this check
	delete(c.assumed, pod.Name)
	return true // re-queued (potential duplicate)
}

var duplicateSchedules int64

func schedule(cache *AssumeCache, pod *Pod, wg *sync.WaitGroup) {
	defer wg.Done()

	err := cache.assume(pod, "node-1")
	if err != nil {
		// Already assumed — this IS the duplicate
		atomic.AddInt64(&duplicateSchedules, 1)
		return
	}

	// Start async bind
	go func() {
		if err := cache.bind(pod); err != nil {
			// Bind failed — trigger error handling
			if cache.handleError(pod) {
				// Pod re-queued — schedule again (this causes the race)
				var innerWg sync.WaitGroup
				innerWg.Add(1)
				go schedule(cache, pod, &innerWg)
				innerWg.Wait()
			}
		}
	}()

	// Simulate concurrent error that also triggers handleError
	// This is the race: bind() and error handler both check Bound
	time.Sleep(time.Millisecond * 3)
	if cache.handleError(pod) {
		// Re-queued from error path while bind is in-flight
		var innerWg sync.WaitGroup
		innerWg.Add(1)
		go schedule(cache, pod, &innerWg)
		innerWg.Wait()
	}
}

func main() {
	cache := NewAssumeCache()
	var wg sync.WaitGroup

	// Schedule 10 pods concurrently — some will hit the race
	for i := 0; i < 10; i++ {
		wg.Add(1)
		pod := &Pod{Name: fmt.Sprintf("pod-%d", i)}
		go schedule(cache, pod, &wg)
	}

	wg.Wait()

	dupes := atomic.LoadInt64(&duplicateSchedules)
	bindings := atomic.LoadInt64(&cache.bindings)

	fmt.Printf("bindings: %d\n", bindings)
	fmt.Printf("duplicate_schedules: %d\n", dupes)

	if dupes > 0 {
		fmt.Println("BUG: duplicate scheduling detected — assume-bind race confirmed")
	} else {
		fmt.Println("OK: no duplicate scheduling")
	}
}
