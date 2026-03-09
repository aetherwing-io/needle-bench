# memory-leak-event-listener

## Project

A Node.js HTTP server that processes data batches through a validation and transformation pipeline. Components communicate via an EventEmitter-based event bus.

## Symptoms

After sending a few hundred requests to POST /process, the /stats endpoint shows that event listener counts grow linearly with the number of requests. Memory (heapUsed) also grows steadily. The server never crashes immediately, but in production it would eventually hit memory limits and be OOM-killed.

## Bug description

Somewhere in the request handling pipeline, event listeners are being registered in a way that accumulates over time. Each request adds to the problem, and nothing ever cleans up the old listeners. The closures retained by these listeners also hold references to request-scoped data, preventing garbage collection.

## Difficulty

Medium

## Expected turns

6-12
