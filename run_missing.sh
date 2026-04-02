#!/usr/bin/env bash
# Run missing benchmarks for a specific model group.
# Usage: bash run_missing.sh <group>
#   Groups: anthropic, google, mistral, openai, xai, deepseek, qwen, other
# Or: bash run_missing.sh --all
#
# SPEC-bench-v2 §8: fire-and-forget — skip-if-exists, log capture, summary.
set -euo pipefail
cd "$(dirname "$0")"

BENCHMARKS=(
  api-version-field-drop
  auth-bypass-path-traversal
  compiler-macro-expansion
  data-corruption-concurrent-write
  deadlock-transfer
  encoding-mojibake
  goroutine-leak-handler
  linearizability-stale-read
  performance-cliff-hash
  raft-snapshot-commit-gap
  silent-data-corruption
  split-brain-leader-election
  sql-injection-search
  timing-attack-comparison
  tls-chain-ordering-strict
  wal-fsync-ghost-ack
)

# Models with native drivers → run native + kernel + kernel-cpu
ANTHROPIC=(claude-haiku-4-5 claude-sonnet-4-6 claude-opus-4-6)
GOOGLE=(gemini-2.5-flash gemini-2.5-pro gemini-3-flash-preview gemini-3.1-pro-preview)
MISTRAL=(codestral-2508 devstral-2512 devstral-medium devstral-small-latest)

# Kernel-only models (OpenRouter) → run kernel only
OPENAI=(gpt-4.1 gpt-5-codex o3 o4-mini)
XAI=(grok-3 grok-3-fast grok-3-mini grok-4 grok-4-fast grok-4.1-fast grok-4.20 grok-code-fast-1)
DEEPSEEK=(deepseek-r1 deepseek-r1-0528 deepseek-v3.2)
QWEN=(qwen3-coder qwen3-coder-flash qwen3-coder-plus)
OTHER=(kimi-k2.5 llama-4-maverick)

GROUP="${1:?Usage: $0 <group|--all>}"

# Counters for summary
TOTAL=0
SKIPPED=0
PASSED=0
FAILED=0

run_bench() {
  local model="$1" bench="$2" arm_label="$3" arm_flags="$4"
  local score_file="runs/${model}-${arm_label}/${bench}.score.json"
  TOTAL=$((TOTAL + 1))

  if [ -f "$score_file" ]; then
    SKIPPED=$((SKIPPED + 1))
    return 0
  fi

  echo "RUN  $model / $bench / $arm_label"
  # shellcheck disable=SC2086
  if ostk bench "$bench" --model "$model" $arm_flags --docker --keep 2>&1; then
    PASSED=$((PASSED + 1))
  else
    FAILED=$((FAILED + 1))
    echo "WARN: failed $model / $bench / $arm_label"
  fi
}

run_driver_models() {
  local -n models_ref=$1
  for model in "${models_ref[@]}"; do
    for bench in "${BENCHMARKS[@]}"; do
      run_bench "$model" "$bench" "native"     "--arm native --local"
      run_bench "$model" "$bench" "kernel"     "--arm kernel --local"
      run_bench "$model" "$bench" "kernel-cpu" "--arm kernel --driver cpu --local"
    done
  done
}

run_kernel_only() {
  local -n models_ref=$1
  for model in "${models_ref[@]}"; do
    for bench in "${BENCHMARKS[@]}"; do
      run_bench "$model" "$bench" "kernel" "--arm kernel --local"
    done
  done
}

START_TIME=$(date +%s)
echo "=== needle-bench fleet: group=$GROUP started=$(date -u +%Y-%m-%dT%H:%M:%SZ) ==="

case "$GROUP" in
  anthropic) run_driver_models ANTHROPIC ;;
  google)    run_driver_models GOOGLE ;;
  mistral)   run_driver_models MISTRAL ;;
  openai)    run_kernel_only OPENAI ;;
  xai)       run_kernel_only XAI ;;
  deepseek)  run_kernel_only DEEPSEEK ;;
  qwen)      run_kernel_only QWEN ;;
  other)     run_kernel_only OTHER ;;
  --all)
    run_driver_models ANTHROPIC
    run_driver_models GOOGLE
    run_driver_models MISTRAL
    run_kernel_only OPENAI
    run_kernel_only XAI
    run_kernel_only DEEPSEEK
    run_kernel_only QWEN
    run_kernel_only OTHER
    ;;
  *) echo "Unknown group: $GROUP"; exit 1 ;;
esac

END_TIME=$(date +%s)
ELAPSED=$((END_TIME - START_TIME))

# Summary
echo ""
echo "=== needle-bench summary ==="
echo "  group:   $GROUP"
echo "  total:   $TOTAL"
echo "  skipped: $SKIPPED (score file exists)"
echo "  ran:     $((PASSED + FAILED))"
echo "  passed:  $PASSED"
echo "  failed:  $FAILED"
echo "  elapsed: ${ELAPSED}s"
echo "  done:    $(date -u +%Y-%m-%dT%H:%M:%SZ)"

# Auto-consolidate scores if any new runs completed
if [ $((PASSED + FAILED)) -gt 0 ] && [ -f consolidate_scores.py ]; then
  echo ""
  echo "=== consolidating scores ==="
  python3 consolidate_scores.py || echo "WARN: consolidation failed"
fi
