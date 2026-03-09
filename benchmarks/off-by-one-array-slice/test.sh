#!/bin/sh
# Test for off-by-one-array-slice bug
# Verifies that batch processing handles all records without dropping any.

set -e

cd /app

FAIL=0

# Test 1: Process 25 records in batches of 10 — should get all 25 back
result=$(python -c "
from processor import process_batch
from data import generate_records
import random
random.seed(42)
records = generate_records(25)
processed = process_batch(records, 10)
print(len(processed))
")

if [ "$result" != "25" ]; then
    echo "FAIL: Expected 25 processed records, got $result"
    FAIL=1
fi

# Test 2: Process 10 records in batches of 10 — exact batch boundary
result2=$(python -c "
from processor import process_batch
from data import generate_records
import random
random.seed(42)
records = generate_records(10)
processed = process_batch(records, 10)
print(len(processed))
")

if [ "$result2" != "10" ]; then
    echo "FAIL: Expected 10 processed records for exact batch size, got $result2"
    FAIL=1
fi

# Test 3: Process 7 records in batches of 3 — uneven division
result3=$(python -c "
from processor import process_batch
from data import generate_records
import random
random.seed(42)
records = generate_records(7)
processed = process_batch(records, 3)
print(len(processed))
")

if [ "$result3" != "7" ]; then
    echo "FAIL: Expected 7 processed records for uneven batches, got $result3"
    FAIL=1
fi

# Test 4: Verify all record IDs are present (no drops)
result4=$(python -c "
from processor import process_batch
from data import generate_records
import random
random.seed(42)
records = generate_records(20)
processed = process_batch(records, 5)
ids = sorted([r['id'] for r in processed])
expected = list(range(1, 21))
print('match' if ids == expected else f'mismatch: got {ids}')
")

if [ "$result4" != "match" ]; then
    echo "FAIL: Record IDs don't match expected sequence: $result4"
    FAIL=1
fi

if [ $FAIL -eq 0 ]; then
    echo "PASS: All records processed correctly across batch boundaries"
    exit 0
else
    exit 1
fi
