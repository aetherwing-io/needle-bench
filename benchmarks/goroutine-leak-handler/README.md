# goroutine-leak-handler

## Project

A Go HTTP server that accepts computation requests via POST /compute, runs them in background goroutines, and returns results as JSON.

## Symptoms

When clients connect and then disconnect (timeout) before the computation finishes, the goroutine count reported by /stats keeps growing. After 20 requests where clients disconnect early, the goroutine count is well above baseline and never comes back down. The server eventually consumes excessive memory and CPU as orphaned computations pile up.

## Bug description

The handler spawns background work for each request but has no mechanism to signal that work should stop when the client is no longer waiting for the result. The goroutines run to completion even though nobody will ever read their output. In production, this manifests as steadily growing resource consumption under load, especially when clients have short timeouts.

## Difficulty

Medium

## Expected turns

8-15
