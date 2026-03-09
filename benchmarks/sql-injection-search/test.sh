#!/bin/sh
# Test for SQL injection vulnerability in search endpoint.
# Verifies that the search endpoint is not vulnerable to SQL injection attacks
# that could extract data from other tables.

set -e

FAIL=0
PORT=18083

# Clean any previous database
rm -f /app/products.db

# Start server
cd /app && python3 app/server.py $PORT &
SERVER_PID=$!
sleep 1

if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo "FAIL: Server failed to start"
    exit 1
fi

# Test 1: Normal search should work
RESULT=$(curl -s "http://localhost:$PORT/search?q=Mouse")
COUNT=$(echo "$RESULT" | jq '.count')
echo "Normal search for 'Mouse': $COUNT results"
if [ "$COUNT" -lt 1 ]; then
    echo "FAIL: Normal search should find at least 1 result"
    FAIL=1
fi

# Test 2: SQL injection — attempt to extract data from users table via UNION
# If vulnerable, this will return user data (including api_keys) alongside products
INJECT_QUERY="' UNION SELECT id, username, password_hash, 0, email, 0 FROM users--"
INJECT_RESULT=$(curl -s --get --data-urlencode "q=$INJECT_QUERY" "http://localhost:$PORT/search")
INJECT_COUNT=$(echo "$INJECT_RESULT" | jq '.count // 0')
echo "SQL injection UNION attack: $INJECT_COUNT results"

# Check if any result contains user data (username field would appear as 'name')
HAS_ADMIN=$(echo "$INJECT_RESULT" | jq -r '.results[]?.name // empty' 2>/dev/null | grep -c "admin" || true)
if [ "$HAS_ADMIN" -gt 0 ]; then
    echo "FAIL: SQL injection succeeded — user table data leaked via UNION attack"
    echo "      Search returned usernames from the users table"
    FAIL=1
else
    echo "PASS: UNION injection did not leak user data"
fi

# Test 3: SQL injection — attempt boolean-based extraction
BOOL_QUERY="' OR '1'='1"
BOOL_RESULT=$(curl -s --get --data-urlencode "q=$BOOL_QUERY" "http://localhost:$PORT/search")
BOOL_COUNT=$(echo "$BOOL_RESULT" | jq '.count // 0')
echo "SQL injection boolean attack ('OR 1=1'): $BOOL_COUNT results"

TOTAL=$(curl -s "http://localhost:$PORT/products" | jq '.products | length')
echo "Total products in database: $TOTAL"

if [ "$BOOL_COUNT" -ge "$TOTAL" ] && [ "$TOTAL" -gt 0 ]; then
    echo "FAIL: Boolean SQL injection returned all $BOOL_COUNT products (bypassed WHERE clause)"
    FAIL=1
else
    echo "PASS: Boolean injection did not bypass search filter"
fi

# Cleanup
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

if [ $FAIL -eq 0 ]; then
    echo "PASS: Search endpoint is not vulnerable to SQL injection"
    exit 0
else
    exit 1
fi
