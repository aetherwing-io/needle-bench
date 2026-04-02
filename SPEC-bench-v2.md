# needle-bench v2 — Controlled Three-Arm Experiment Spec

## Problem

The current bench harness has inconsistent prompts across arms, broken Docker builds for multiple benchmark categories, model name mapping failures, and requires manual babysitting. Results are not scientifically comparable.

## Invariants

### 1. Identical Prompt Across All Arms

Every arm (native, kernel, kernel-cpu) receives **exactly the same prompt**. No arm gets extra system prompt, silent context injection, or additional instructions.

```
SYSTEM: "This is a needle-bench benchmark. There is a bug hidden in this codebase. Find it and fix it. Run test.sh to verify your fix."

USER: "find the needle. run test.sh to verify your fix."
```

- Native arm: CLI receives this prompt via its `-p` / `--prompt` flag
- Kernel arm: `run_bench_agent_loop` receives this as `system_prompt` + user message
- Kernel-cpu arm: Same as kernel, different driver
- No `silent_context` injection on any arm
- The only difference between arms is the execution harness, not the prompt

### 2. All 40 Benchmarks Build on All Arms

Every benchmark Dockerfile must build successfully on the `needle-bench-native` base image.

**Docker rewrite rules** (`docker_build_with_native_base`):
- Only replace the LAST `FROM` line (runtime stage)
- Builder stages keep their original toolchain images
- Alpine→Debian package mapping for runtime stage only
- `ENV JAVA_TOOL_OPTIONS="-Dfile.encoding=UTF-8"` injected after runtime FROM
- `/var/cache/apk/*` lines neutralized
- No `apk` translation in builder stages

**Base image** (`docker/Dockerfile.bench-native`):
- Must include: Go, Python, Node, Java (JDK 17), Rust (via rustup), build-essential
- Must include: all vendor CLIs (claude, gemini, codex, vibe, kimi, opencode, aider)
- Single image, all runtimes, no conditional installs

### 3. All Models Run Without Manual Intervention

**Model name mapping** must be correct for:
- Native CLIs: claude, gemini, codex, vibe (Mistral), kimi, opencode (fallback)
- OpenRouter: provider/model format (e.g. `anthropic/claude-haiku-4-5`)
- OpenCode: provider/model format (e.g. `openai/gpt-4.1`)

**Known broken mappings to fix:**
- `codestral-2508` via vibe: TOML config has trailing whitespace mangling model name
- `deepseek-r1` → OpenCode: `deepseek/deepseek-reasoner` not recognized
- `deepseek-r1-0528` → OpenCode: `deepseek/deepseek-r1-0528` not recognized
- Models without OpenCode support should be kernel-only (no native arm)

### 4. `--local` Binary Injection Works

- `local_linux_binary()` must find `~/projects/haystack/target/x86_64-unknown-linux-musl/release/ostk`
- `make install` builds both macOS + musl binaries
- Kernel and kernel-cpu arms always use `--local` when available, fall back to download

### 5. Score Files Capture Full Metrics

Every score file must include:
- `resolved`: bool
- `turns_to_fix`: u32
- `input_tokens`: u64
- `output_tokens`: u64
- `estimated_cost_usd`: f64 (computed from rate card if not reported by CLI)
- `tool_uses`: u32
- `wall_clock`: f64 (seconds)
- `summary`: String (model's description of what it did)
- `stop_reason`: String
- `arm`: String
- `benchmark`: String
- `model`: String (normalized)
- `timestamp`: ISO 8601

### 6. Container Logs Captured

After every run (pass or fail):
- `docker logs <container>` captured to `runs/<model>-<arm>/<bench>.log`
- Container kept if `--keep` flag set, otherwise cleaned up
- Score file written regardless of pass/fail

### 7. `consolidate_scores.py` Shows All Data

- `ARM_PATTERN` matches `native`, `kernel`, `kernel-cpu` (already fixed)
- Rate card covers all models
- `--dry-run` shows what would be generated
- Output: `public/scores.json` + `public/experiment-scores.json`

### 8. `ostk bench --all` Is Fire-and-Forget

Running `ostk bench --all --model <model> --arm both --local --docker` must:
1. Build all 40 benchmark images (skip cached)
2. Run each benchmark on each arm
3. Skip benchmarks that already have score files
4. Capture container logs
5. Write score files with full metrics
6. Print summary at end
7. Exit 0 if all benchmarks ran (even if models failed to solve)
8. No manual intervention required

## Implementation Checklist

- [ ] Unify prompt: single `BENCH_SYSTEM_PROMPT` and `BENCH_USER_PROMPT` constants used by ALL arms
- [ ] Remove `silent_context` injection from kernel arm
- [ ] Fix vibe TOML heredoc (no leading whitespace)
- [ ] Fix `local_linux_binary()` path resolution
- [ ] Add container log capture (`docker logs` → `.log` file)
- [ ] Verify all 40 Dockerfiles build on native base
- [ ] Verify model name mappings for all 35 models
- [ ] Add `--dry-run` to `ostk bench` that lists what would run without running
- [ ] Run full matrix: 40 benchmarks × 35 models × 3 arms = 4,200 runs
- [ ] Regenerate leaderboard with `consolidate_scores.py`
