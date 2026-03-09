"""
CLI entry point for the scheduler. Loads events from a JSON config
and reports which events are due at a given time.
"""

import json
import sys
from datetime import datetime

import pytz

from models import ScheduledEvent
from event_store import EventStore
from scheduler import find_due_events


def load_events(config_path: str) -> EventStore:
    """Load events from a JSON configuration file."""
    store = EventStore()

    with open(config_path, "r") as f:
        data = json.load(f)

    for entry in data.get("events", []):
        event = ScheduledEvent(
            event_id=entry["id"],
            name=entry["name"],
            hour=entry["hour"],
            minute=entry["minute"],
            timezone=entry["timezone"],
            owner=entry.get("owner", "system"),
            enabled=entry.get("enabled", True),
        )
        store.add_event(event)

    return store


def main():
    if len(sys.argv) < 3:
        print("Usage: python runner.py <config.json> <utc_datetime>")
        print("  utc_datetime format: YYYY-MM-DDTHH:MM")
        sys.exit(1)

    config_path = sys.argv[1]
    time_str = sys.argv[2]

    try:
        now_utc = datetime.strptime(time_str, "%Y-%m-%dT%H:%M")
        now_utc = pytz.utc.localize(now_utc)
    except ValueError:
        print(f"ERROR: Invalid datetime format: {time_str}")
        sys.exit(2)

    store = load_events(config_path)
    total = len(store.list_enabled())
    due = find_due_events(store, now_utc, window_minutes=1)

    print(f"Checking {total} events at {now_utc.strftime('%Y-%m-%d %H:%M UTC')}")
    print(f"Due events: {len(due)}")

    for event in due:
        print(f"  FIRE: {event}")

    if not due:
        print("  (no events due)")


if __name__ == "__main__":
    main()
