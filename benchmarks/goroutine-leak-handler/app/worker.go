package main

import (
	"math"
	"sync"
)

// WorkerPool manages a pool of reusable computation workers.
type WorkerPool struct {
	mu      sync.Mutex
	workers []*Worker
	maxSize int
}

// Worker represents a single computation unit.
type Worker struct {
	ID       int
	busy     bool
	computed int64
}

// NewWorkerPool creates a pool with the given max size.
func NewWorkerPool(maxSize int) *WorkerPool {
	return &WorkerPool{
		workers: make([]*Worker, 0, maxSize),
		maxSize: maxSize,
	}
}

// Acquire gets an available worker or creates one if under capacity.
func (p *WorkerPool) Acquire() *Worker {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, w := range p.workers {
		if !w.busy {
			w.busy = true
			return w
		}
	}

	if len(p.workers) < p.maxSize {
		w := &Worker{ID: len(p.workers), busy: true}
		p.workers = append(p.workers, w)
		return w
	}

	return nil
}

// Release returns a worker to the pool.
func (p *WorkerPool) Release(w *Worker) {
	p.mu.Lock()
	defer p.mu.Unlock()
	w.busy = false
}

// Compute runs a computation on the worker. This is pure math, no I/O.
func (w *Worker) Compute(iterations int) float64 {
	sum := 0.0
	for i := 0; i < iterations; i++ {
		sum += math.Sqrt(float64(i+1)) * math.Log(float64(i+2))
		w.computed++
	}
	return sum
}
