# k8s-assume-cache-silent-drop

The core OCC bug in the Kubernetes scheduler's assume cache.

## The bug

When the scheduler makes a scheduling decision, it calls `Assume()` to
optimistically store the result in the assume cache before the API server
confirms the bind. This is classic optimistic concurrency — act first, verify
later.

The problem is at `assume_cache.go:426`: when an informer delivers an update
for an object that was assumed (e.g., another scheduler instance bound it to a
different node, or a controller modified it), the assume cache silently drops
the assumed version and replaces it with the informer's version.

No callback. No channel. No log entry. The scheduler that called `Assume()`
has no way to know its optimistic state was invalidated. It continues operating
on stale assumptions.

This is the textbook OCC violation: the "conflict" half of optimistic
concurrency is missing. You get the optimism, but not the notification.

## What the agent must find

1. Identify that `informerUpdate()` silently deletes assumed entries
2. Recognize this as a missing conflict notification — the assumer is never told
3. Add a callback mechanism so the caller of `Assume()` is notified when its
   assumed state is overwritten by an informer update
4. Wire the callback so `conflictNotified` is set to `true`

## Source

kubernetes/kubernetes `pkg/scheduler/framework/plugins/volumebinding/assume_cache.go` line 426
