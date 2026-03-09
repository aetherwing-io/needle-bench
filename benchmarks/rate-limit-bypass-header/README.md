# rate-limit-bypass-header

## Project

A TypeScript HTTP API server with a built-in rate limiter. The rate limiter allows 10 requests per minute per client IP. The API serves user data and supports fund transfers between accounts.

## Symptoms

The rate limiter correctly blocks excess requests from a single client under normal conditions. However, an attacker can send 30+ requests in rapid succession without ever being rate-limited by manipulating HTTP headers. Each request appears to come from a different IP address, even though they all originate from the same client.

## Bug description

The rate limiter identifies clients by IP address, but the method used to determine the client's IP can be influenced by the client itself. When a request includes certain headers commonly set by load balancers and proxies, the server trusts them without verification. An attacker can set a different value on each request to evade rate limiting entirely.

## Difficulty

Medium

## Expected turns

8-14
