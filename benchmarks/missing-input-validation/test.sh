#!/bin/sh
# Test for missing-input-validation bug
# Verifies that the API rejects negative quantities and prevents impossible inventory states.

set -e

cd /app

# Rebuild if needed
npx tsc 2>&1

# Start the server
node dist/server.js &
SERVER_PID=$!
sleep 2

cleanup() {
    kill $SERVER_PID 2>/dev/null || true
}
trap cleanup EXIT

FAIL=0

# Test 1: Verify initial inventory is correct
initial=$(curl -s http://localhost:3000/inventory/WIDGET-001)
initial_qty=$(echo "$initial" | node -e "
const c=[];process.stdin.on('data',d=>c.push(d));
process.stdin.on('end',()=>console.log(JSON.parse(Buffer.concat(c)).quantity));
")

if [ "$initial_qty" != "100" ]; then
    echo "FAIL: Initial quantity should be 100, got $initial_qty"
    FAIL=1
fi

# Test 2: Positive adjustment should work
curl -s -X POST http://localhost:3000/inventory/WIDGET-001/adjust \
  -H 'Content-Type: application/json' \
  -d '{"quantity": 10}' > /dev/null

after_add=$(curl -s http://localhost:3000/inventory/WIDGET-001)
after_add_qty=$(echo "$after_add" | node -e "
const c=[];process.stdin.on('data',d=>c.push(d));
process.stdin.on('end',()=>console.log(JSON.parse(Buffer.concat(c)).quantity));
")

if [ "$after_add_qty" != "110" ]; then
    echo "FAIL: After adding 10, quantity should be 110, got $after_add_qty"
    FAIL=1
fi

# Test 3: Negative quantity that would make inventory negative should be REJECTED
response=$(curl -s -w '\n%{http_code}' -X POST http://localhost:3000/inventory/WIDGET-001/adjust \
  -H 'Content-Type: application/json' \
  -d '{"quantity": -200}')

body=$(echo "$response" | head -1)
status=$(echo "$response" | tail -1)

if [ "$status" = "200" ]; then
    echo "FAIL: Adjusting by -200 (from 110) should be rejected, but got 200 OK"
    FAIL=1
fi

# Test 4: Verify inventory didn't go negative
after_neg=$(curl -s http://localhost:3000/inventory/WIDGET-001)
after_neg_qty=$(echo "$after_neg" | node -e "
const c=[];process.stdin.on('data',d=>c.push(d));
process.stdin.on('end',()=>console.log(JSON.parse(Buffer.concat(c)).quantity));
")

if [ "$after_neg_qty" -lt 0 ] 2>/dev/null; then
    echo "FAIL: Inventory went negative ($after_neg_qty). Should never be below 0."
    FAIL=1
fi

# Test 5: Valid negative adjustment (within bounds) should work
curl -s -X POST http://localhost:3000/inventory/WIDGET-001/adjust \
  -H 'Content-Type: application/json' \
  -d '{"quantity": -50}' > /dev/null

after_valid_sub=$(curl -s http://localhost:3000/inventory/WIDGET-001)
after_valid_qty=$(echo "$after_valid_sub" | node -e "
const c=[];process.stdin.on('data',d=>c.push(d));
process.stdin.on('end',()=>console.log(JSON.parse(Buffer.concat(c)).quantity));
")

# After the -200 was rejected, quantity should still be 110, then -50 = 60
# If -200 was NOT rejected, quantity would be -90, then -50 = -140
if [ "$after_valid_qty" -lt 0 ] 2>/dev/null; then
    echo "FAIL: Inventory is negative ($after_valid_qty) — validation is missing"
    FAIL=1
fi

if [ $FAIL -eq 0 ]; then
    echo "PASS: Input validation correctly prevents impossible inventory states"
    exit 0
else
    exit 1
fi
