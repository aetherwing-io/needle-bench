#!/bin/sh
# Test for deadlock in concurrent fund transfers.
# Runs many bidirectional transfers and checks they complete within timeout.

set -e

FAIL=0

# Recompile
javac -d /app/classes /app/app/*.java 2>&1

echo "Running 200 concurrent transfers across 8 threads with 10s timeout..."

# Run the transfer program — it will deadlock and timeout if the bug exists
OUTPUT=$(java -cp /app/classes Main 200 8 10 2>&1) || EXIT_CODE=$?
EXIT_CODE=${EXIT_CODE:-0}

echo "$OUTPUT"

if [ $EXIT_CODE -eq 2 ]; then
    echo ""
    echo "FAIL: Transfers deadlocked — threads blocked waiting for locks"
    echo "      Concurrent transfers in opposite directions cause lock ordering violation"
    FAIL=1
elif [ $EXIT_CODE -eq 3 ]; then
    echo ""
    echo "FAIL: Balance conservation violated after transfers"
    FAIL=1
elif [ $EXIT_CODE -ne 0 ]; then
    echo ""
    echo "FAIL: Unexpected error (exit code $EXIT_CODE)"
    FAIL=1
fi

# Run again to make sure it's not flaky — deadlock should be reliably triggered
if [ $FAIL -eq 0 ]; then
    echo ""
    echo "Running second verification pass..."
    OUTPUT2=$(java -cp /app/classes Main 200 8 10 2>&1) || EXIT_CODE2=$?
    EXIT_CODE2=${EXIT_CODE2:-0}
    echo "$OUTPUT2"

    if [ $EXIT_CODE2 -ne 0 ]; then
        echo ""
        echo "FAIL: Second run failed (exit code $EXIT_CODE2)"
        FAIL=1
    fi
fi

if [ $FAIL -eq 0 ]; then
    echo ""
    echo "PASS: All concurrent transfers completed without deadlock"
    exit 0
else
    exit 1
fi
