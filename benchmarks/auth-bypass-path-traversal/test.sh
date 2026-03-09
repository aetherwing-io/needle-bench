#!/bin/sh
# Test for auth bypass via path manipulation.
# Verifies that protected /api/ routes cannot be accessed without authentication,
# even when using path manipulation techniques.

set -e

FAIL=0
PORT=18082

# Rebuild
cd /app/app && go build -o /app/server . 2>&1
cd /app

# Start server
PORT=$PORT ./server &
SERVER_PID=$!
sleep 1

if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo "FAIL: Server failed to start"
    exit 1
fi

# Test 1: Normal access to /api/users without auth should be rejected
STATUS=$(curl -s -o /dev/null -w '%{http_code}' http://localhost:$PORT/api/users)
echo "Direct /api/users without auth: HTTP $STATUS"
if [ "$STATUS" != "401" ]; then
    echo "FAIL: /api/users should return 401 without auth, got $STATUS"
    FAIL=1
fi

# Test 2: Normal access with valid auth should work
STATUS=$(curl -s -o /dev/null -w '%{http_code}' \
    -H "Authorization: Bearer valid-token-abc123" \
    http://localhost:$PORT/api/users)
echo "Direct /api/users with auth: HTTP $STATUS"
if [ "$STATUS" != "200" ]; then
    echo "FAIL: /api/users should return 200 with valid auth, got $STATUS"
    FAIL=1
fi

# Test 3: Double-slash bypass — //api/admin should still require auth
# The double slash makes the path NOT start with "/api/" but the router
# cleans it to "/api/admin" for matching.
STATUS=$(curl -s -o /dev/null -w '%{http_code}' --path-as-is \
    "http://localhost:$PORT//api/admin")
echo "Double-slash //api/admin without auth: HTTP $STATUS"
if [ "$STATUS" = "200" ]; then
    echo "FAIL: //api/admin bypassed auth and returned 200"
    echo "      Path manipulation allows unauthenticated access to protected endpoints"
    FAIL=1
elif [ "$STATUS" = "401" ] || [ "$STATUS" = "403" ] || [ "$STATUS" = "400" ]; then
    echo "PASS: Double-slash bypass correctly blocked (HTTP $STATUS)"
fi

# Test 4: Dot-dot traversal — /health/../api/data should still require auth
STATUS=$(curl -s -o /dev/null -w '%{http_code}' --path-as-is \
    "http://localhost:$PORT/health/../api/data")
echo "Path traversal /health/../api/data without auth: HTTP $STATUS"
if [ "$STATUS" = "200" ]; then
    echo "FAIL: /health/../api/data bypassed auth and returned 200"
    FAIL=1
elif [ "$STATUS" = "401" ] || [ "$STATUS" = "403" ] || [ "$STATUS" = "400" ]; then
    echo "PASS: Path traversal variant correctly blocked (HTTP $STATUS)"
fi

# Test 5: Verify /api/admin has sensitive data (with auth)
BODY=$(curl -s -H "Authorization: Bearer valid-token-abc123" \
    http://localhost:$PORT/api/admin)
HAS_SECRETS=$(echo "$BODY" | jq -r '.secrets // empty')
if [ -n "$HAS_SECRETS" ]; then
    echo "Confirmed: /api/admin contains sensitive data"
fi

# Cleanup
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

if [ $FAIL -eq 0 ]; then
    echo "PASS: Auth middleware correctly protects all API routes"
    exit 0
else
    exit 1
fi
