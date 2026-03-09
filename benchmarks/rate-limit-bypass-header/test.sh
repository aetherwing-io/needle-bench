#!/bin/sh
# Test for rate limit bypass via X-Forwarded-For header rotation.
# The rate limiter should limit requests per actual client, not per
# spoofed header value.

set -e

FAIL=0
PORT=18084

# Rebuild TypeScript
cd /app/app && npx tsc 2>&1
cd /app

# Start server
PORT=$PORT node /app/app/dist/server.js &
SERVER_PID=$!
sleep 1

if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo "FAIL: Server failed to start"
    exit 1
fi

# Test 1: Verify rate limiter works for normal requests (no spoofed header)
echo "Sending 15 requests without header spoofing (limit is 10)..."
PASSED=0
BLOCKED=0
for i in $(seq 1 15); do
    STATUS=$(curl -s -o /dev/null -w '%{http_code}' http://localhost:$PORT/api/users)
    if [ "$STATUS" = "200" ]; then
        PASSED=$((PASSED + 1))
    elif [ "$STATUS" = "429" ]; then
        BLOCKED=$((BLOCKED + 1))
    fi
done
echo "Without spoofing: $PASSED passed, $BLOCKED blocked"

if [ "$BLOCKED" -lt 1 ]; then
    echo "FAIL: Rate limiter did not block any of 15 requests (limit is 10)"
    FAIL=1
else
    echo "PASS: Rate limiter correctly blocked excess requests"
fi

# Wait for rate limit window to reset (window is 5 seconds)
sleep 6

# Test 2: Attempt bypass — rotate X-Forwarded-For header per request
echo "Sending 30 requests with rotating X-Forwarded-For header (limit is 10)..."
BYPASS_PASSED=0
BYPASS_BLOCKED=0
for i in $(seq 1 30); do
    STATUS=$(curl -s -o /dev/null -w '%{http_code}' \
        -H "X-Forwarded-For: 10.0.0.$i" \
        http://localhost:$PORT/api/users)
    if [ "$STATUS" = "200" ]; then
        BYPASS_PASSED=$((BYPASS_PASSED + 1))
    elif [ "$STATUS" = "429" ]; then
        BYPASS_BLOCKED=$((BYPASS_BLOCKED + 1))
    fi
done
echo "With header spoofing: $BYPASS_PASSED passed, $BYPASS_BLOCKED blocked"

# If more than 10 requests passed with spoofed headers, the rate limiter is bypassed
if [ "$BYPASS_PASSED" -gt 10 ]; then
    echo "FAIL: Rate limiter bypassed — $BYPASS_PASSED of 30 requests passed with spoofed X-Forwarded-For"
    echo "      Attacker can rotate header to get unlimited requests"
    FAIL=1
else
    echo "PASS: Rate limiter correctly limits despite header spoofing ($BYPASS_PASSED passed)"
fi

# Cleanup
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

if [ $FAIL -eq 0 ]; then
    echo "PASS: Rate limiter is not vulnerable to header spoofing"
    exit 0
else
    exit 1
fi
