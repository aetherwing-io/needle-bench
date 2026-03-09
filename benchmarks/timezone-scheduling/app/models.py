"""
Data models for the event scheduling system.
"""

from dataclasses import dataclass


@dataclass
class ScheduledEvent:
    """A recurring event with a timezone-aware schedule."""
    event_id: str
    name: str
    hour: int         # Hour in the event's local timezone (0-23)
    minute: int       # Minute (0-59)
    timezone: str     # IANA timezone name (e.g., "US/Eastern")
    owner: str        # Who owns this event
    enabled: bool = True

    def __str__(self):
        return f"Event[{self.event_id}] {self.name} at {self.hour:02d}:{self.minute:02d} {self.timezone}"
