"""
In-memory event store. Manages the list of scheduled events.
"""

from models import ScheduledEvent


class EventStore:
    """Simple in-memory store for scheduled events."""

    def __init__(self):
        self._events = {}

    def add_event(self, event: ScheduledEvent):
        """Add or update a scheduled event."""
        self._events[event.event_id] = event

    def remove_event(self, event_id: str):
        """Remove an event by ID."""
        self._events.pop(event_id, None)

    def get_event(self, event_id: str) -> ScheduledEvent | None:
        """Retrieve a single event."""
        return self._events.get(event_id)

    def list_events(self) -> list[ScheduledEvent]:
        """Return all events."""
        return list(self._events.values())

    def list_enabled(self) -> list[ScheduledEvent]:
        """Return only enabled events."""
        return [e for e in self._events.values() if e.enabled]
