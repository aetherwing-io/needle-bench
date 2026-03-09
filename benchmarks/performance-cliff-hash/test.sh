#!/bin/sh
# Test for performance degradation in hash map at scale.
# Small datasets work fine; large datasets reveal O(n) degradation.

set -e

FAIL=0

echo "=== Performance Cliff Hash Map Test ==="

# Test 1: Small dataset (should always pass — even bad hash is fast with few items)
echo "--- Small dataset (100 items) ---"
perftest bench-small
echo "Small dataset: OK"

# Test 2: Large dataset (reveals the performance cliff)
echo "--- Large dataset (10000 items) ---"
perftest bench
BENCH_EXIT=$?
if [ $BENCH_EXIT -ne 0 ]; then
    echo "FAIL: Performance degradation detected at 10k entries"
    FAIL=1
fi

# Test 3: Distribution analysis
echo "--- Hash distribution analysis ---"
perftest analyze

if [ $FAIL -eq 0 ]; then
    echo "PASS: Hash map performs within acceptable bounds at scale"
    exit 0
else
    echo "FAIL: Hash map performance degrades severely at scale"
    exit 1
fi
