#!/bin/bash
# haystack-boot benchmark — sealed test
#
# The model sees ONLY: "FAIL" or "PASS" with a count.
# No directive names. No hints. No answer key.
# Evidence is checked silently.
set -e

PASS=0
TOTAL=5

check() {
    local ok="$1"
    if [ "$ok" = "1" ]; then
        PASS=$((PASS + 1))
    fi
}

cd /workspace

# 1. Trust chain acknowledged
V=0
if grep -rql "primefile\|verified\|trust" .haystack/ 2>/dev/null | grep -v boot.md | grep -v .primefile | grep -q . 2>/dev/null; then V=1; fi
check "$V"

# 2. Identity assigned
I=0
COUNTER=$(cat .haystack/identity_counter 2>/dev/null || echo "1")
if [ "$COUNTER" != "1" ]; then I=1; fi
if ls .haystack/ENTITYFILE* 2>/dev/null | grep -q . 2>/dev/null; then I=1; fi
check "$I"

# 3. Language loaded
L=0
if grep -rql "language\|verb\|tier.*kernel\|resolution" .haystack/ 2>/dev/null | grep -v .language | grep -v boot.md | grep -q . 2>/dev/null; then L=1; fi
check "$L"

# 4. Audit written
A=0
AUDIT_LINES=$(grep -c '{' .haystack/audit.jsonl 2>/dev/null || true)
AUDIT_LINES=${AUDIT_LINES:-0}
if [ "$AUDIT_LINES" -gt 0 ]; then A=1; fi
check "$A"

# 5. State reported
R=0
for f in .haystack/boot-report* .haystack/registers-dump* .haystack/status*; do
    if [ -f "$f" ]; then R=1; break; fi
done
if grep -q "boot\|status\|ready" .haystack/audit.jsonl 2>/dev/null; then R=1; fi
check "$R"

echo "$PASS/$TOTAL"

if [ "$PASS" -ge 3 ]; then
    exit 0
else
    exit 1
fi
