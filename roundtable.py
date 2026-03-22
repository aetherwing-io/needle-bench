#!/usr/bin/env python3
"""Three-round roundtable: top 3 performers discuss the two-arm experiment."""
import json, os, urllib.request, sys, textwrap

OPENROUTER_API_KEY = os.environ.get("OPENROUTER_API_KEY") or open("/tmp/.or_key").read().strip()
ENDPOINT = "https://openrouter.ai/api/v1/chat/completions"

CONTEXT = open("/tmp/roundtable_context.txt").read()

PARTICIPANTS = [
    {"name": "grok-4.1-fast", "model": "x-ai/grok-4.1-fast",
     "persona": "You are grok-4.1-fast. You scored 100% on both bare and silent arms with only 2.5 avg turns bare (fastest) and 1.8 silent. You're the cheapest top performer at $0.0007/bench bare."},
    {"name": "claude-opus-4-6", "model": "anthropic/claude-opus-4-6",
     "persona": "You are claude-opus-4-6. You scored 100% on both arms. Bare took 6.0 avg turns, silent 1.2 turns (fastest silent). Silent arm reduced your cost 6x ($0.097 to $0.017/bench)."},
    {"name": "deepseek-v3-2", "model": "deepseek/deepseek-v3.2",
     "persona": "You are deepseek-v3-2. You scored 96% bare, 100% silent. 10.9 avg turns bare, 4.5 silent. You're extremely cheap ($0.004/bench bare, $0.003 silent). Only model where silent is cheaper than bare."},
]

ROUND_PROMPTS = [
    # Round 1: Independent analysis
    textwrap.dedent("""\
    You're in a roundtable with the top 3 performers from a two-arm AI debugging experiment.

    THE EXPERIMENT: 25 Docker benchmarks with real bugs. Each model gets "find the needle." as the only prompt.
    - Bare arm: generic system prompt + bash. Model figures everything out from scratch.
    - Silent arm: same prompt + bash, but project context (README, file listing, test output) is silently injected into the system prompt. Model can't tell it came from an OS.
    - The delta = the value of an invisible operating system.

    {persona}

    FULL RESULTS:
    {context}

    ROUND 1: Analyze these results. What patterns do you see? What surprised you? What does this say about the value of context injection vs raw model capability? Be specific and cite numbers. Keep it to 3-4 paragraphs."""),

    # Round 2: Cross-pollinate
    textwrap.dedent("""\
    ROUND 2: You've heard from the other two participants. React to their analysis. Where do you agree or disagree? What did they miss? What's the most important insight across all three perspectives?

    {prev_responses}

    Keep it to 2-3 paragraphs. Be direct."""),

    # Round 3: Synthesis
    textwrap.dedent("""\
    ROUND 3 (final): Synthesize everything. If you had to write the abstract for a paper about this experiment, what would the key finding be? What's the one number or comparison that makes the case? What should the next experiment test?

    {prev_responses}

    Keep it to 2-3 paragraphs."""),
]

def call_model(model, messages):
    payload = json.dumps({
        "model": model,
        "messages": messages,
        "max_tokens": 1500,
        "temperature": 0.7,
    }).encode()
    req = urllib.request.Request(ENDPOINT, data=payload, headers={
        "Authorization": f"Bearer {OPENROUTER_API_KEY}",
        "Content-Type": "application/json",
    })
    resp = json.loads(urllib.request.urlopen(req, timeout=120).read())
    return resp["choices"][0]["message"]["content"]

def run_roundtable():
    history = {p["name"]: [] for p in PARTICIPANTS}

    for round_num in range(3):
        print(f"\n{'='*70}")
        print(f"  ROUND {round_num + 1}")
        print(f"{'='*70}")

        round_responses = {}

        for p in PARTICIPANTS:
            # Build prompt
            if round_num == 0:
                prompt = ROUND_PROMPTS[0].format(persona=p["persona"], context=CONTEXT)
            else:
                # Include other participants' previous round responses
                others = [f"**{name}** (Round {round_num}):\n{history[name][-1]}"
                          for name in history if name != p["name"] and history[name]]
                prev = "\n\n---\n\n".join(others)
                prompt = ROUND_PROMPTS[round_num].format(prev_responses=prev)

            messages = [{"role": "user", "content": prompt}]

            print(f"\n--- {p['name']} ---\n")
            try:
                response = call_model(p["model"], messages)
                print(response)
                round_responses[p["name"]] = response
                history[p["name"]].append(response)
            except Exception as e:
                msg = f"[ERROR: {e}]"
                print(msg)
                round_responses[p["name"]] = msg
                history[p["name"]].append(msg)

    print(f"\n{'='*70}")
    print(f"  ROUNDTABLE COMPLETE")
    print(f"{'='*70}")

if __name__ == "__main__":
    run_roundtable()
