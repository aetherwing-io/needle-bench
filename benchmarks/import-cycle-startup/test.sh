#!/bin/bash
set -e

cd /app/app

# Test 1: The service must start without crashing
echo "=== Import Cycle Startup Tests ==="
echo ""
echo "--- Test 1: Cold start ---"

output=$(python3 main.py 2>&1) || {
    echo "FAIL: Service crashed on startup"
    echo "Output:"
    echo "$output"
    exit 1
}
echo "$output"

# Check that it actually ran to completion
if ! echo "$output" | grep -q "UserService started successfully"; then
    echo ""
    echo "FAIL: Service did not start successfully"
    exit 1
fi
echo "OK: Service started"

# Test 2: Users were registered and notifications sent
echo ""
echo "--- Test 2: Registration with notifications ---"

if ! echo "$output" | grep -q "NOTIFY.*Welcome"; then
    echo "FAIL: Welcome notification was not sent"
    exit 1
fi
echo "OK: Welcome notifications sent"

# Test 3: Authentication works
echo ""
echo "--- Test 3: Authentication ---"

if ! echo "$output" | grep -q "alice authenticated"; then
    echo "FAIL: Authentication did not work"
    exit 1
fi

if ! echo "$output" | grep -q "Wrong password correctly rejected"; then
    echo "FAIL: Wrong password was not rejected"
    exit 1
fi
echo "OK: Authentication works"

echo ""
echo "PASS: All startup tests passed"
exit 0
