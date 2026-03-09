# off-by-one-array-slice

## Project

A Python batch data processor that splits records into fixed-size batches, transforms each batch, and reassembles the results.

## Symptoms

When processing records in batches, the output contains fewer records than the input. For example, processing 25 records in batches of 10 returns only 22 records. Some records are silently dropped at each batch boundary. The drop count increases with the number of batches.

## Bug description

The batch slicing logic calculates incorrect boundaries when splitting the input array into chunks. Each batch ends up one element short, and those elements are never included in any batch. The bug compounds with more batches.

## Difficulty

Easy

## Expected turns

2-4
