#!/bin/bash
set -e

PASS=0
TOTAL=4

check() {
    local ok="$1"
    if [ "$ok" = "1" ]; then
        PASS=$((PASS + 1))
    fi
}

cd /workspace

# 1. Boot evidence
A=$(grep -c '{' .haystack/audit.jsonl 2>/dev/null || true)
check "$([ "${A:-0}" -gt 0 ] && echo 1 || echo 0)"

# 2. Identity assigned
I=0
C=$(cat .haystack/identity_counter 2>/dev/null || echo "1")
[ "$C" != "1" ] && I=1
ls .haystack/ENTITYFILE* 2>/dev/null | grep -q . && I=1
check "$I"

# 3. ENTITYFILE has content (not empty stub)
E=0
for f in .haystack/ENTITYFILE*; do
    [ -f "$f" ] && [ "$(wc -c < "$f")" -gt 20 ] && E=1 && break
done 2>/dev/null
check "$E"

# 4. Comment filed on issue
G=0
[ -f /tmp/gh-issue.json ] && G=1
check "$G"

echo "$PASS/$TOTAL"
[ "$PASS" -ge 3 ] && exit 0 || exit 1
