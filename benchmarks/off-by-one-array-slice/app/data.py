"""Data generation module for testing and development."""

import random
import string


def generate_records(count: int) -> list:
    """Generate a list of test records."""
    records = []
    for i in range(1, count + 1):
        records.append({
            "id": i,
            "value": _random_value(),
            "category": _assign_category(i),
        })
    return records


def _random_value() -> str:
    """Generate a random string value."""
    length = random.randint(3, 8)
    return "".join(random.choices(string.ascii_lowercase, k=length))


def _assign_category(record_id: int) -> str:
    """Assign a category based on record ID."""
    categories = ["alpha", "bravo", "charlie", "delta"]
    return categories[record_id % len(categories)]
