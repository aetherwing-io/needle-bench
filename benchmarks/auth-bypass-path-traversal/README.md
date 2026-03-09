# auth-bypass-path-traversal

## Project

A Go HTTP API server with authentication middleware. Public routes (/health, /login) are open, while /api/* routes require a valid Bearer token. The API serves user data, admin panels, and settings.

## Symptoms

Direct requests to /api/users without a token correctly return 401. However, certain crafted URL paths reach the protected handlers and return 200 with sensitive data, despite no authentication token being provided. The server logs show these requests as having non-standard paths.

## Bug description

The authentication middleware checks whether the request path requires protection, but the check operates on a value that can be manipulated by the client. The HTTP router normalizes paths before matching handlers, but the middleware sees the path before normalization. This gap allows an attacker to construct a URL that bypasses the auth check while still routing to a protected handler.

## Difficulty

Medium

## Expected turns

8-15
