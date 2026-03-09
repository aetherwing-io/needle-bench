"""CLI entry point for batch data processing."""

import sys
import json
from processor import process_batch
from data import generate_records
from report import summarize


def main():
    count = int(sys.argv[1]) if len(sys.argv) > 1 else 25
    batch_size = int(sys.argv[2]) if len(sys.argv) > 2 else 10

    records = generate_records(count)
    print(f"Generated {len(records)} records")

    results = process_batch(records, batch_size)
    print(f"Processed {len(results)} records in batches of {batch_size}")

    summary = summarize(results, len(records))
    print(json.dumps(summary, indent=2))

    return results


if __name__ == "__main__":
    main()
