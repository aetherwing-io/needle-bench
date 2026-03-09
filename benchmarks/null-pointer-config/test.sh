#!/bin/sh
# Test for null-pointer-config bug
# The server should start and respond correctly when metrics are enabled in config.

set -e

cd /app

# Build the server
go build -o server . 2>&1

# Start the server in the background
./server &
SERVER_PID=$!
sleep 2

cleanup() {
    kill $SERVER_PID 2>/dev/null || true
}
trap cleanup EXIT

FAIL=0

# Test 1: health endpoint should work
health=$(curl -s -o /dev/null -w '%{http_code}' http://localhost:8080/health)
if [ "$health" != "200" ]; then
    echo "FAIL: /health returned $health, expected 200"
    FAIL=1
fi

# Test 2: data endpoint should work
data_code=$(curl -s -o /dev/null -w '%{http_code}' http://localhost:8080/data)
if [ "$data_code" != "200" ]; then
    echo "FAIL: /data returned $data_code, expected 200"
    FAIL=1
fi

# Test 3: status endpoint should return metrics info when enabled
status_code=$(curl -s -o /dev/null -w '%{http_code}' http://localhost:8080/status)
if [ "$status_code" != "200" ]; then
    echo "FAIL: /status returned $status_code, expected 200"
    FAIL=1
else
    status_body=$(curl -s http://localhost:8080/status)
    has_metrics=$(echo "$status_body" | grep -c '"metrics"' || true)
    if [ "$has_metrics" -lt 1 ]; then
        echo "FAIL: /status response missing metrics field when metrics are enabled"
        FAIL=1
    fi
fi

if [ $FAIL -eq 0 ]; then
    echo "PASS: Server starts and responds correctly with metrics enabled"
    exit 0
else
    exit 1
fi
