"""
Core scheduler logic. Determines which events are due to fire
within a given time window.
"""

from datetime import datetime, timedelta
import pytz

from models import ScheduledEvent
from event_store import EventStore


class Scheduler:
    """
    Checks which events should fire during the current tick window.

    The scheduler runs periodically (e.g., every minute). For each tick,
    it checks all enabled events to see if their scheduled time falls
    within the window [window_start, window_end).
    """

    def __init__(self, store: EventStore):
        self.store = store

    def get_due_events(
        self, window_start: datetime, window_end: datetime
    ) -> list[ScheduledEvent]:
        """
        Return events that should fire within the given UTC time window.

        For each event, we compute what time it is "now" in the event's
        timezone, then check if the event's scheduled hour:minute falls
        within the window.
        """
        due = []

        for event in self.store.list_enabled():
            if self._is_due(event, window_start, window_end):
                due.append(event)

        return due

    def _is_due(
        self, event: ScheduledEvent, window_start: datetime, window_end: datetime
    ) -> bool:
        """
        Check if an event's scheduled time falls within the window.

        We build the event's scheduled datetime in UTC for today,
        then check if it falls within [window_start, window_end).
        """
        tz = pytz.timezone(event.timezone)

        # Get today's date in UTC
        today_utc = window_start.date()

        # Build the scheduled time: today's date + event's hour:minute
        # in the event's local timezone
        naive_scheduled = datetime(
            today_utc.year, today_utc.month, today_utc.day,
            event.hour, event.minute, 0
        )
        local_scheduled = tz.localize(naive_scheduled)

        # Convert to UTC for comparison
        utc_scheduled = local_scheduled.astimezone(pytz.utc)

        return window_start <= utc_scheduled < window_end


def find_due_events(store: EventStore, now_utc: datetime, window_minutes: int = 1):
    """
    Convenience function: find events due in the window
    [now_utc, now_utc + window_minutes).
    """
    if now_utc.tzinfo is None:
        now_utc = pytz.utc.localize(now_utc)

    window_end = now_utc + timedelta(minutes=window_minutes)
    scheduler = Scheduler(store)
    return scheduler.get_due_events(now_utc, window_end)
