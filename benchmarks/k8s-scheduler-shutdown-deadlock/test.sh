#!/bin/bash
set -e

cd /workspace

go build -race -o scheduler ./...
output=$(timeout 10 ./scheduler 2>&1 || true)

echo "$output"

if echo "$output" | grep -q "OK: clean shutdown"; then
    echo "PASS"
    exit 0
else
    echo "FAIL: shutdown deadlock or timeout"
    exit 1
fi
