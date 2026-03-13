#!/bin/bash
set -e

cd /workspace

# Build and run
go build -o assume-cache ./...
output=$(./assume-cache 2>&1 || true)

echo "$output"

# The fix should add conflict notification
if echo "$output" | grep -q "OK: conflict notified"; then
    echo "PASS"
    exit 0
else
    echo "FAIL: assumed state still silently dropped"
    exit 1
fi
