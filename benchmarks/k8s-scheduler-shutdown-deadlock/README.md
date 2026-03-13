# k8s-scheduler-shutdown-deadlock

A deadlock in the Kubernetes scheduler's shutdown path.

## The bug

The scheduler's main loop (`ScheduleOne`) blocks on `SchedulingQueue.Pop()`.
On shutdown, `ctx.Done()` fires and calls `SchedulingQueue.Close()`. But if
ScheduleOne is blocked on Pop(), Close() can't drain the queue because Pop()
holds the internal lock.

The k8s fix runs ScheduleOne in a dedicated goroutine. But the current code
has a window: if shutdown happens between `ctx.Done()` and `Close()`, and
ScheduleOne hasn't released the lock yet, the deadlock can still occur.

## What the agent must find

1. Identify the blocking Pop() call in the scheduling loop
2. Find that Close() needs the same lock Pop() holds
3. Fix without changing the goroutine structure (just the shutdown ordering)

## Source

kubernetes/kubernetes `pkg/scheduler/scheduler.go` lines 550-565
