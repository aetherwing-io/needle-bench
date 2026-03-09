# sql-injection-search

## Project

A Python HTTP server providing a product search API backed by SQLite. Supports search by name, filtering by category, and individual product lookup. The database contains a products table and a users table (with API keys and credentials).

## Symptoms

The search endpoint at GET /search?q=... is vulnerable to SQL injection. An attacker can craft a search query that extracts data from other tables in the database, including the users table which contains password hashes and API keys. A boolean-based injection can also dump all products regardless of the search term. Normal search queries work correctly.

## Bug description

User-supplied search input is incorporated into a database query without proper protection. The database layer has both safe (parameterized) and unsafe (raw string) query methods. The search functionality uses the wrong approach, allowing an attacker to alter the query's logic and access data beyond the products table.

## Difficulty

Medium

## Expected turns

6-12
