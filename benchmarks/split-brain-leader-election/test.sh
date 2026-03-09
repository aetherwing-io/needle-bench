#!/bin/sh
# Test for split-brain in leader election under network partition.
# A correct implementation must never elect two leaders simultaneously.

set -e

FAIL=0

echo "=== Split-Brain Leader Election Test ==="

# Run the safety check multiple times — the simulation is probabilistic
for round in 1 2 3; do
    echo "--- Round $round ---"
    if ! split-brain check; then
        echo "FAIL: Split-brain detected in round $round"
        FAIL=1
        break
    fi
    echo "Round $round: OK"
done

if [ $FAIL -eq 0 ]; then
    echo "PASS: No split-brain detected across all rounds"
    exit 0
else
    echo "FAIL: Leader election allows split-brain under partition"
    exit 1
fi
