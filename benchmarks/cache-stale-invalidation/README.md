# cache-stale-invalidation

## Project

A Python HTTP server implementing a product catalog API with an in-memory caching layer. Supports CRUD operations on products. The cache uses a 5-minute TTL to reduce repeated computation.

## Symptoms

After updating a product via PUT /products/:id, a subsequent GET /products/:id returns the old data. After creating a new product via POST /products, a subsequent GET /products does not include it. The changes are persisted in the store, but reads continue to return outdated values. The problem resolves itself after about 5 minutes.

## Bug description

The caching layer correctly caches reads and has a TTL, but the write path does not coordinate with the cache. When data is modified through the API, the cached version becomes stale but is still served to subsequent readers. The cache has the methods to handle this, but the write endpoints do not invoke them.

## Difficulty

Medium

## Expected turns

6-10
