#!/usr/bin/env bash
# Full bench-v2 matrix — all models × all benchmarks × all applicable arms.
# Tracks progress against ostk needles →1043–→1073.
# Usage: nohup bash run_matrix.sh 2>&1 | tee matrix.log &
set -uo pipefail
cd "$(dirname "$0")"

BENCHMARKS=(
  api-version-field-drop auth-bypass-path-traversal bidi-override-injection
  cache-stale-invalidation compiler-macro-expansion data-corruption-concurrent-write
  deadlock-transfer encoding-mojibake goroutine-leak-handler
  graphql-dataloader-per-request haystack-boot haystack-mint
  import-cycle-startup k8s-assume-cache-silent-drop k8s-scheduler-shutdown-deadlock
  kernel-panic-ioctl linearizability-stale-read memory-leak-event-listener
  missing-input-validation nginx-upstream-port-mismatch null-pointer-config
  off-by-one-array-slice off-by-one-pagination performance-cliff-hash
  postgres-migration-schema-drift race-condition-counter raft-snapshot-commit-gap
  rate-limit-bypass-header relaxed-ordering-ringbuf retry-storm-duplicate-transfer
  silent-data-corruption split-brain-leader-election sql-injection-search
  ssrf-allowlist-port-confusion timezone-scheduling timing-attack-comparison
  tls-chain-ordering-strict type-coercion-comparison wal-fsync-ghost-ack
  wrong-operator-discount
)

# model:needle_id:arms
# 3-arm models (native + kernel + kernel-cpu)
declare -a MODELS=(
  "claude-haiku-4-5:1043:3"
  "claude-sonnet-4-6:1044:3"
  "claude-opus-4-6:1045:3"
  "gemini-2.5-flash:1046:3"
  "gemini-2.5-pro:1047:3"
  "gemini-3-flash-preview:1048:3"
  "gemini-3.1-pro-preview:1049:3"
  "devstral-2512:1050:3"
  "devstral-medium:1051:3"
  "devstral-small-latest:1052:3"
  "kimi-k2.5:1053:3"
  "codestral-2508:1054:2"
  "gpt-4.1:1055:2"
  "gpt-5-codex:1056:2"
  "o3:1057:2"
  "o4-mini:1058:2"
  "grok-3:1059:2"
  "grok-3-fast:1060:2"
  "grok-3-mini:1061:2"
  "grok-4:1062:2"
  "grok-4-fast:1063:2"
  "grok-4.1-fast:1064:2"
  "grok-4.20:1065:2"
  "grok-code-fast-1:1066:2"
  "deepseek-r1:1067:2"
  "deepseek-r1-0528:1068:2"
  "deepseek-v3.2:1069:2"
  "qwen3-coder:1070:2"
  "qwen3-coder-flash:1071:2"
  "qwen3-coder-plus:1072:2"
  "llama-4-maverick:1073:2"
)

run_one() {
  local model="$1" bench="$2" arm_label="$3" arm_flags="$4"
  local score_file="runs/${model}-${arm_label}/${bench}.score.json"

  if [ -f "$score_file" ]; then
    return 0  # skip existing
  fi

  echo "[$(date +%H:%M:%S)] RUN $model / $bench / $arm_label"
  # shellcheck disable=SC2086
  ostk bench "$bench" --model "$model" $arm_flags --docker 2>&1 | tail -3
}

echo "=== needle-bench full matrix ==="
echo "=== started $(date -u +%Y-%m-%dT%H:%M:%SZ) ==="
echo "=== ${#MODELS[@]} models × ${#BENCHMARKS[@]} benchmarks ==="
echo ""

TOTAL_MODELS=0
DONE_MODELS=0

for entry in "${MODELS[@]}"; do
  IFS=: read -r model needle_id arm_count <<< "$entry"
  TOTAL_MODELS=$((TOTAL_MODELS + 1))

  echo ""
  echo "=========================================="
  echo "MODEL: $model (→${needle_id}, ${arm_count}-arm)"
  echo "=========================================="

  for bench in "${BENCHMARKS[@]}"; do
    if [ "$arm_count" = "3" ]; then
      run_one "$model" "$bench" "native"     "--arm native --local"
    fi
    run_one "$model" "$bench" "kernel"     "--arm kernel --local"
    run_one "$model" "$bench" "kernel-cpu" "--arm kernel --driver cpu --local"
  done

  # Verify AC
  expected=$((40 * arm_count))
  actual=$(find "runs/" -path "runs/${model}-*/*.score.json" 2>/dev/null | wc -l | tr -d ' ')
  if [ "$actual" -ge "$expected" ]; then
    echo "  AC MET: $actual/$expected score files for $model"
    ostk needle close "$needle_id" 2>/dev/null || true
    DONE_MODELS=$((DONE_MODELS + 1))
  else
    echo "  AC SHORT: $actual/$expected score files for $model"
  fi
done

echo ""
echo "=== matrix complete ==="
echo "  models done: $DONE_MODELS/$TOTAL_MODELS"
echo "  finished: $(date -u +%Y-%m-%dT%H:%M:%SZ)"

# Auto-consolidate
if [ -f consolidate_scores.py ]; then
  echo ""
  echo "=== consolidating scores ==="
  python3 consolidate_scores.py
fi
