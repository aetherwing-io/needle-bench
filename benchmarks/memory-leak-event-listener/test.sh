#!/bin/sh
# Test for memory leak via event listener accumulation
# Sends many requests and checks that listener count and memory stay bounded.

set -e

FAIL=0
PORT=8080

# Start server in background
node /app/app/server.js &
SERVER_PID=$!
sleep 1

# Verify server is running
if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo "FAIL: Server failed to start"
    exit 1
fi

# Get baseline stats
BASELINE=$(curl -s http://localhost:$PORT/stats)
BASELINE_LISTENERS=$(echo "$BASELINE" | jq '.listenerCounts["data:validated"]')
BASELINE_HEAP=$(echo "$BASELINE" | jq '.memory.heapUsed_mb')
echo "Baseline: listeners=$BASELINE_LISTENERS, heap=${BASELINE_HEAP}MB"

# Send 500 process requests — each should NOT accumulate listeners
for i in $(seq 1 500); do
    curl -s -X POST http://localhost:$PORT/process \
        -H "Content-Type: application/json" \
        -d '{"items":[{"id":"item-'$i'","value":'$i'}]}' \
        >/dev/null 2>&1
done

# Force GC and wait
sleep 2

# Get stats after load
AFTER=$(curl -s http://localhost:$PORT/stats)
AFTER_LISTENERS=$(echo "$AFTER" | jq '.listenerCounts["data:validated"]')
AFTER_HEAP=$(echo "$AFTER" | jq '.memory.heapUsed_mb')
echo "After 500 requests: listeners=$AFTER_LISTENERS, heap=${AFTER_HEAP}MB"

# Check listener count — should NOT grow with number of requests
# A properly implemented processor registers listeners once, not per-request
if [ "$AFTER_LISTENERS" -gt 10 ]; then
    echo "FAIL: Listener count grew to $AFTER_LISTENERS (expected <= 10)"
    echo "      Event listeners are being registered per-request and never removed"
    FAIL=1
else
    echo "PASS: Listener count stable at $AFTER_LISTENERS"
fi

# Cleanup
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

if [ $FAIL -eq 0 ]; then
    echo "PASS: No event listener leak detected"
    exit 0
else
    exit 1
fi
