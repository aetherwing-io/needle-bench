# k8s-scheduler-assume-race

A race condition in the Kubernetes scheduler's assume-and-bind path.

## The bug

When a pod is assumed (optimistically placed on a node before binding), the
scheduler continues scheduling other pods. If the assume fails during retry,
there's a window where the pod exists in the assume cache but not in the
actual binding — creating a phantom placement that blocks real scheduling.

The code at `schedule_one.go:327` documents this as "most probably result of
a BUG in retrying logic." The fix relies on Error() checking if the pod was
already bound, but this check races with the binding goroutine.

## What the agent must find

1. Identify the race window between assume() and bind()
2. Find that Error() can re-queue an already-bound pod
3. Fix the race without adding locks (k8s uses optimistic concurrency)

## Source

kubernetes/kubernetes `pkg/scheduler/schedule_one.go` lines 320-335
