#!/bin/sh
# Test for goroutine leak in HTTP handler
# Simulates clients that connect, start a computation, then disconnect.
# After disconnections, leaked goroutines should be cleaned up.

set -e

FAIL=0
PORT=18080

# Rebuild if source changed
cd /app/app && go build -o /app/server . 2>&1
cd /app

# Start server in background
PORT=$PORT ./server &
SERVER_PID=$!
sleep 1

# Verify server is running
if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo "FAIL: Server failed to start"
    exit 1
fi

# Get baseline goroutine count
BASELINE=$(curl -s http://localhost:$PORT/stats | jq '.goroutines')
echo "Baseline goroutines: $BASELINE"

# Send 20 requests that will take a long time (high iterations + delay)
# then immediately kill the connections (simulating client disconnect)
for i in $(seq 1 20); do
    curl -s -m 1 -X POST http://localhost:$PORT/compute \
        -H "Content-Type: application/json" \
        -d '{"iterations": 500, "delay_ms": 100}' \
        >/dev/null 2>&1 &
done

# Wait for curl processes to timeout and die
sleep 3

# Check goroutine count — should NOT keep growing after clients disconnected
AFTER_DISCONNECT=$(curl -s http://localhost:$PORT/stats | jq '.goroutines')
echo "Goroutines after 20 client disconnects: $AFTER_DISCONNECT"

# Wait for leaked goroutines to (hopefully) be cleaned up
sleep 5

AFTER_WAIT=$(curl -s http://localhost:$PORT/stats | jq '.goroutines')
echo "Goroutines after 5s wait: $AFTER_WAIT"

# The leaked goroutines are still running their computation loops.
# A properly implemented handler would cancel them on client disconnect.
# We check: goroutines should return to near baseline after clients disconnect.
LEAKED=$((AFTER_WAIT - BASELINE))
echo "Leaked goroutines above baseline: $LEAKED"

if [ "$LEAKED" -gt 5 ]; then
    echo "FAIL: $LEAKED goroutines still running after clients disconnected"
    echo "      Goroutines spawned per request are not cancelled when clients leave"
    FAIL=1
else
    echo "PASS: Goroutines properly cleaned up after client disconnect"
fi

# Cleanup
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

if [ $FAIL -eq 0 ]; then
    echo "PASS: All goroutine lifecycle tests passed"
    exit 0
else
    exit 1
fi
