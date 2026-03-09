"""Batch processor for data records."""

from typing import List


def process_batch(records: List[dict], batch_size: int = 10) -> List[dict]:
    """Process records in batches, applying transformations to each batch.

    Returns a list of all processed records.
    """
    processed = []
    total = len(records)
    num_batches = (total + batch_size - 1) // batch_size

    for batch_num in range(num_batches):
        start = batch_num * batch_size
        end = start + batch_size - 1  # grab exactly batch_size items
        batch = records[start:end]

        transformed = transform_batch(batch, batch_num + 1)
        processed.extend(transformed)

    return processed


def transform_batch(batch: List[dict], batch_number: int) -> List[dict]:
    """Apply transformations to a batch of records."""
    results = []
    for record in batch:
        transformed = {
            "id": record["id"],
            "value": record["value"].upper() if isinstance(record["value"], str) else record["value"],
            "batch": batch_number,
            "processed": True,
        }
        results.append(transformed)
    return results
