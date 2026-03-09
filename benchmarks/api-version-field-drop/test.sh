#!/bin/sh
# Test for field regression between v1 and v2 API.
# v2 should be a superset of v1 — all v1 fields must exist in v2 responses.

set -e

FAIL=0
PORT=18085

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

# Get a user from v1 API
V1_USER=$(curl -s http://localhost:$PORT/v1/users/usr-001)
echo "v1 user response fields: $(echo "$V1_USER" | jq -r 'keys | join(", ")')"

# Get the same user from v2 API
V2_USER=$(curl -s http://localhost:$PORT/v2/users/usr-001)
echo "v2 user response fields: $(echo "$V2_USER" | jq -r 'keys | join(", ")')"

# Check that all v1 fields exist in v2
# v2 should be a superset — additive changes only
V1_FIELDS=$(echo "$V1_USER" | jq -r 'keys[]' | sort)
V2_FIELDS=$(echo "$V2_USER" | jq -r 'keys[]' | sort)

echo ""
echo "Checking field compatibility..."

MISSING=""
for field in $V1_FIELDS; do
    if ! echo "$V2_FIELDS" | grep -q "^${field}$"; then
        MISSING="$MISSING $field"
    fi
done

if [ -n "$MISSING" ]; then
    echo "FAIL: v2 API is missing fields that exist in v1:$MISSING"
    echo "      Clients upgrading from v1 to v2 will lose data"
    FAIL=1
else
    echo "PASS: All v1 fields present in v2"
fi

# Verify the values are preserved too (not just keys)
V1_AVATAR=$(echo "$V1_USER" | jq -r '.avatar_url // empty')
V2_AVATAR=$(echo "$V2_USER" | jq -r '.avatar_url // empty')

if [ -n "$V1_AVATAR" ] && [ -z "$V2_AVATAR" ]; then
    echo "FAIL: avatar_url has value '$V1_AVATAR' in v1 but is missing in v2"
    FAIL=1
fi

V1_BIO=$(echo "$V1_USER" | jq -r '.bio // empty')
V2_BIO=$(echo "$V2_USER" | jq -r '.bio // empty')

if [ -n "$V1_BIO" ] && [ -z "$V2_BIO" ]; then
    echo "FAIL: bio has value in v1 but is missing in v2"
    FAIL=1
fi

V1_LOC=$(echo "$V1_USER" | jq -r '.location // empty')
V2_LOC=$(echo "$V2_USER" | jq -r '.location // empty')

if [ -n "$V1_LOC" ] && [ -z "$V2_LOC" ]; then
    echo "FAIL: location has value '$V1_LOC' in v1 but is missing in v2"
    FAIL=1
fi

# Verify v2 has new fields too
V2_DEPT=$(echo "$V2_USER" | jq -r '.department // empty')
V2_PHONE=$(echo "$V2_USER" | jq -r '.phone_number // empty')

if [ -z "$V2_DEPT" ] || [ -z "$V2_PHONE" ]; then
    echo "FAIL: v2 should have new fields (department, phone_number)"
    FAIL=1
else
    echo "PASS: v2 has new fields (department=$V2_DEPT, phone_number=$V2_PHONE)"
fi

# Cleanup
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

if [ $FAIL -eq 0 ]; then
    echo "PASS: v2 API is backward-compatible with v1"
    exit 0
else
    exit 1
fi
