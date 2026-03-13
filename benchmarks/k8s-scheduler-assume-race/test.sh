#!/bin/bash
set -e

cd /workspace

# Build and run with race detector
go build -race -o scheduler ./...
output=$(./scheduler 2>&1 || true)

echo "$output"

# The fix should eliminate duplicate scheduling
if echo "$output" | grep -q "OK: no duplicate scheduling"; then
    echo "PASS"
    exit 0
else
    echo "FAIL: race condition still present"
    exit 1
fi
