#!/bin/sh
# Test for data corruption in concurrent file writer.
# Writes a file using multiple concurrent threads and verifies
# each segment contains the expected deterministic fill pattern.

set -e

FAIL=0
OUTPUT=/tmp/concurrent_test.dat

echo "=== Concurrent Write Integrity Test ==="

# Run the concurrent write 5 times — race conditions are probabilistic
for run in 1 2 3 4 5; do
    rm -f "$OUTPUT"
    concurrent-writer write "$OUTPUT"

    if ! concurrent-writer verify "$OUTPUT"; then
        echo "FAIL: Corruption detected on run $run"
        FAIL=1
        break
    fi
    echo "Run $run: OK"
done

rm -f "$OUTPUT"

if [ $FAIL -eq 0 ]; then
    echo "PASS: All concurrent write runs produced correct output"
    exit 0
else
    echo "FAIL: Data corruption from concurrent writes"
    exit 1
fi
