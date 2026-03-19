# needle-bench

**Your worst debugging day, everyone's benchmark.**

A benchmark suite of 29 scenarios for AI coding agents, built from real bugs in real codebases. No synthetic tasks. No contrived puzzles. Just broken containers and one prompt: *find the needle.*

## How it works

Each benchmark is a Docker container with a real bug. The agent gets tools (`shell`, `file:read`, `file:edit`), a time limit, and a test that fails. The agent explores, diagnoses, and patches. The test either passes or it doesn't.

```
benchmarks/off-by-one-pagination/
  Dockerfile              # broken codebase, containerized
  Agentfile               # agent config: tools, limits
  .bench/solution.patch   # sealed truth (agent never sees this)
  test.sh                 # exit 0 = fixed, exit 1 = broken
```

Scenarios span concurrency bugs, security bypasses, encoding issues, k8s operational failures, and more.

## Quick start

```bash
# Validate all benchmarks
make validate

# Run a specific benchmark
make run BENCH=off-by-one-pagination

# List available benchmarks
make list
```

## 11 metrics, no opinions

Every run produces the same 11 numbers. See [SCORING.md](docs/SCORING.md).

| Metric | What it measures |
|--------|-----------------|
| resolved | Did the agent fix the bug? |
| turns_to_discovery | How fast did it find the right file? |
| turns_to_fix | How fast did it produce a working patch? |
| signal_to_noise | What fraction of actions were productive? |
| false_positives | How many wrong files did it edit? |
| token_cost | Total tokens consumed |
| tokens_per_correct_line | Efficiency per correct change |
| recovery_events | How many times did it self-correct? |
| recovery_rate | How often did self-correction succeed? |
| wall_clock | Total time |
| blind_discovery | Did it find the bug with no hints? |

## Submit a benchmark

Your worst debugging day is everyone's benchmark. See [CONTRIBUTING.md](docs/CONTRIBUTING.md).

```bash
cp -r benchmarks/_template benchmarks/your-bug-name
# Edit the files, then:
make validate BENCH=your-bug-name
```

## Leaderboard

Scores are published at [needle-bench.cc](https://needle-bench.cc). 24 models evaluated across 28 benchmarks so far.

Primary rank: resolve rate. Tiebreaker: fewer turns, then fewer tokens.

To regenerate the public leaderboard from individual run scores:

```bash
python3 consolidate_scores.py          # writes public/scores.json
python3 consolidate_scores.py --dry-run # preview without writing
```

## Spec

The full benchmark format specification is in [SPEC.md](docs/SPEC.md).

## License

Apache 2.0. See [LICENSE](LICENSE).

---

*[os-tack/find-the-needle](https://github.com/os-tack/find-the-needle) -- built by Claude Code.*
