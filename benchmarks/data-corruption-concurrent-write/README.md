# data-corruption-concurrent-write

## Project

A Rust CLI tool (`concurrent-writer`) that writes data to a file using multiple worker threads. Each thread is responsible for writing deterministic byte patterns to assigned segments of the file. Used in a data pipeline where parallel I/O is critical for throughput.

## Symptoms

When the tool writes a file and then verifies it, certain segments contain unexpected byte values. The corruption is not deterministic — sometimes the file verifies correctly, sometimes specific segments show wrong fill bytes. The corrupted segments tend to be near boundaries between worker assignments. Running the write + verify cycle multiple times reliably reproduces the issue.

## Bug description

The concurrent write architecture assigns overlapping work to threads without proper synchronization. Multiple threads race to write the same byte ranges, and depending on scheduling, one thread's writes can overwrite another's final state. The corruption appears at segment boundaries where worker responsibilities overlap. Understanding the segment assignment logic and the thread interaction is required to identify the root cause.

## Difficulty

Hard

## Expected turns

10-15
