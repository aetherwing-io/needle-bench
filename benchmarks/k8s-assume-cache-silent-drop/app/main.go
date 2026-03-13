package main

import (
	"fmt"
	"sync"
)

// Simplified reproduction of the k8s assume_cache.go:426 bug.
// When an informer delivers an update while an object is assumed,
// the assumed version gets silently dropped. No notification.
// No conflict. Just data loss.

type Object struct {
	Key     string
	Version int
	Data    string
}

type AssumeCache struct {
	mu      sync.Mutex
	store   map[string]*Object // ground truth from informer
	assumed map[string]*Object // optimistically assumed state
}

func NewAssumeCache() *AssumeCache {
	return &AssumeCache{
		store:   make(map[string]*Object),
		assumed: make(map[string]*Object),
	}
}

// Assume optimistically stores a version of the object.
// The scheduler calls this after making a scheduling decision,
// before the API server confirms the bind.
func (c *AssumeCache) Assume(obj *Object) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.assumed[obj.Key] = obj
	return nil
}

// Get returns the assumed version if it exists, otherwise the store version.
func (c *AssumeCache) Get(key string) (*Object, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if obj, ok := c.assumed[key]; ok {
		return obj, true
	}
	if obj, ok := c.store[key]; ok {
		return obj, true
	}
	return nil, false
}

// informerUpdate is called when the informer delivers a newer version
// of the object from the API server. This is the bug site.
//
// Real k8s code at assume_cache.go:426:
//   If the informer delivers an update for an object that was assumed,
//   the code checks if the informer version is "newer" and if so,
//   removes the assumed entry and replaces it with the informer version.
//   But it does this SILENTLY — no callback, no channel, no log.
//   The scheduler never learns that its assumed state was dropped.
func (c *AssumeCache) informerUpdate(obj *Object) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// BUG: When informer delivers an update for an assumed object,
	// we silently drop the assumed version. The caller who called
	// Assume() is never notified that their optimistic state is gone.
	if assumed, exists := c.assumed[obj.Key]; exists {
		if obj.Version > assumed.Version {
			// Silent drop — no notification to the assumer
			delete(c.assumed, obj.Key)
		}
	}

	c.store[obj.Key] = obj
}

// IsAssumed checks if an object is still in the assumed state.
func (c *AssumeCache) IsAssumed(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.assumed[key]
	return ok
}

var conflictNotified bool

func main() {
	cache := NewAssumeCache()

	// Step 1: Object exists in store at version 3
	cache.informerUpdate(&Object{
		Key:     "pv-0",
		Version: 3,
		Data:    "original",
	})

	// Step 2: Scheduler assumes a new version (e.g., after a scheduling decision)
	// This is the optimistic write — scheduler believes pv-0 is now at version 5
	err := cache.Assume(&Object{
		Key:     "pv-0",
		Version: 5,
		Data:    "assumed-by-scheduler",
	})
	if err != nil {
		fmt.Printf("assume failed: %v\n", err)
		return
	}

	// Verify assumed state is visible
	obj, _ := cache.Get("pv-0")
	fmt.Printf("after assume: version=%d data=%q\n", obj.Version, obj.Data)

	// Step 3: Informer delivers version 6 from the API server.
	// This could be from another scheduler, a controller, or a user edit.
	// The assume cache sees version 6 > assumed version 5,
	// silently drops the assumed entry.
	cache.informerUpdate(&Object{
		Key:     "pv-0",
		Version: 6,
		Data:    "informer-update",
	})

	// Step 4: Check what happened
	obj, _ = cache.Get("pv-0")
	fmt.Printf("after informer: version=%d data=%q\n", obj.Version, obj.Data)

	wasAssumed := cache.IsAssumed("pv-0")
	fmt.Printf("still assumed: %v\n", wasAssumed)

	// The bug: the scheduler's assumed state was silently dropped.
	// It has no way to know. No callback. No channel. No error.
	// It will continue operating as if pv-0 is in its assumed state,
	// but the cache now returns the informer's version instead.
	if !wasAssumed && !conflictNotified {
		fmt.Println("BUG: assumed state silently dropped")
	} else if conflictNotified {
		fmt.Println("OK: conflict notified")
	} else {
		fmt.Println("BUG: unexpected state")
	}
}
