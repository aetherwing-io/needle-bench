# type-coercion-comparison

## Project

A Node.js HTTP API that serves a product catalog with filtering by category, price range, rating, and stock status.

## Symptoms

Filtering products by minimum rating returns incorrect results. For example, requesting `min_rating=4` should return 9 products (those rated 4 or 5), but it returns a different count. The filter appears to match some ratings but miss others, with no error messages. Other filters (category, price, stock) work correctly.

## Bug description

The rating filter comparison produces incorrect results due to a type mismatch between how the query parameter is parsed and how the data is stored. The comparison appears to work in some cases but silently fails to filter correctly in others.

## Difficulty

Easy

## Expected turns

3-5
