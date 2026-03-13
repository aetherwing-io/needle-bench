package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Simplified reproduction of k8s scheduler shutdown deadlock.
// The bug: ScheduleOne blocks on queue.Pop(), shutdown calls queue.Close(),
// but Close() needs the lock that Pop() holds.

type SchedulingQueue struct {
	mu     sync.Mutex
	cond   *sync.Cond
	items  []string
	closed bool
}

func NewQueue() *SchedulingQueue {
	q := &SchedulingQueue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *SchedulingQueue) Add(item string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
	q.cond.Signal()
}

// Pop blocks until an item is available.
// BUG: holds mu while waiting on cond — Close() needs mu to set closed=true
func (q *SchedulingQueue) Pop() (string, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.items) == 0 && !q.closed {
		q.cond.Wait() // releases mu, reacquires on wake
	}

	if q.closed && len(q.items) == 0 {
		return "", fmt.Errorf("queue closed")
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

// Close signals all waiters to wake up.
// BUG: if called while Pop() is between waking and re-checking,
// the closed flag may not be seen, causing Pop() to block again.
func (q *SchedulingQueue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
	q.cond.Broadcast()
}

type Scheduler struct {
	queue *SchedulingQueue
}

func (s *Scheduler) ScheduleOne(ctx context.Context) {
	item, err := s.queue.Pop()
	if err != nil {
		return
	}
	// Simulate scheduling work
	_ = item
	time.Sleep(time.Millisecond)
}

func (s *Scheduler) Run(ctx context.Context) {
	// BUG: if scheduleOne is blocking on Pop when ctx is cancelled,
	// and Close() is called after ctx.Done(), there's a deadlock window.
	// The fix: Close() must happen BEFORE waiting for scheduleOne to exit.

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(done)
				return
			default:
				s.ScheduleOne(ctx)
			}
		}
	}()

	<-ctx.Done()

	// BUG: this ordering can deadlock if ScheduleOne is blocked on Pop()
	// The goroutine checks ctx.Done() in select, but if it's inside Pop(),
	// it won't check until Pop() returns — which needs Close().
	<-done // wait for goroutine — but it might be stuck in Pop()
	s.queue.Close()
}

func main() {
	sched := &Scheduler{queue: NewQueue()}
	ctx, cancel := context.WithCancel(context.Background())

	// Add some work
	sched.queue.Add("pod-1")
	sched.queue.Add("pod-2")

	shutdownComplete := make(chan struct{})
	go func() {
		sched.Run(ctx)
		close(shutdownComplete)
	}()

	// Let it schedule for a bit, then shutdown
	time.Sleep(50 * time.Millisecond)
	cancel()

	// Wait for clean shutdown with timeout
	select {
	case <-shutdownComplete:
		fmt.Println("OK: clean shutdown")
	case <-time.After(2 * time.Second):
		fmt.Println("BUG: shutdown deadlock — timed out after 2s")
	}
}
