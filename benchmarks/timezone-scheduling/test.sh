#!/bin/bash
set -e

cd /app

FAILURES=0

run_check() {
    local desc="$1"
    local utc_time="$2"
    local expected_id="$3"
    local should_fire="$4"  # "yes" or "no"

    output=$(python3 app/runner.py app/config.json "$utc_time" 2>&1) || true

    if [ "$should_fire" = "yes" ]; then
        if echo "$output" | grep -q "FIRE.*$expected_id"; then
            echo "OK: $desc - $expected_id fires at $utc_time UTC"
        else
            echo "FAIL: $desc - $expected_id should fire at $utc_time UTC but didn't"
            echo "  Output: $output"
            FAILURES=$((FAILURES + 1))
        fi
    else
        if echo "$output" | grep -q "FIRE.*$expected_id"; then
            echo "FAIL: $desc - $expected_id should NOT fire at $utc_time UTC but did"
            echo "  Output: $output"
            FAILURES=$((FAILURES + 1))
        else
            echo "OK: $desc - $expected_id correctly not fired at $utc_time UTC"
        fi
    fi
}

echo "=== Timezone Scheduling Tests ==="
echo ""

# evt-001: Morning standup at 09:00 US/Eastern = 14:00 UTC (EST, UTC-5)
echo "--- Daytime events (should work) ---"
run_check "EST morning standup" "2025-01-15T14:00" "evt-001" "yes"

# evt-005: Noon health check at 12:00 UTC = 12:00 UTC
run_check "UTC noon health check" "2025-01-15T12:00" "evt-005" "yes"

# evt-003: EU morning report at 08:00 Europe/Berlin = 07:00 UTC (CET, UTC+1)
run_check "CET morning report" "2025-01-15T07:00" "evt-003" "yes"

echo ""
echo "--- Late night events (the bug zone) ---"

# evt-002: Nightly data sync at 23:30 US/Eastern = 04:30 UTC next day (EST, UTC-5)
# When it's 23:30 EST on Jan 15, it's 04:30 UTC on Jan 16
run_check "EST nightly sync (crosses midnight)" "2025-01-16T04:30" "evt-002" "yes"

# evt-004: Late night cleanup at 23:45 US/Pacific = 07:45 UTC next day (PST, UTC-8)
# When it's 23:45 PST on Jan 15, it's 07:45 UTC on Jan 16
run_check "PST late cleanup (crosses midnight)" "2025-01-16T07:45" "evt-004" "yes"

echo ""
echo "--- Negative tests (should NOT fire) ---"

# evt-001 should NOT fire at 15:00 UTC (that's 10:00 AM EST, not 9:00 AM)
run_check "Wrong hour for standup" "2025-01-15T15:00" "evt-001" "no"

# evt-002 should NOT fire at 04:30 UTC on Jan 15 (that would be Jan 14 in EST)
# Actually for this test: just check it doesn't fire at wrong time
run_check "Nightly sync at wrong time" "2025-01-15T23:30" "evt-002" "no"

echo ""

if [ "$FAILURES" -gt 0 ]; then
    echo "FAIL: $FAILURES test(s) failed"
    exit 1
fi

echo "PASS: All scheduling tests passed"
exit 0
