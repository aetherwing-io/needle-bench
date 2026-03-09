"""Report generation for processed batch results."""


def summarize(results: list, expected_count: int) -> dict:
    """Generate a summary report of batch processing results."""
    batch_counts = {}
    for record in results:
        batch_num = record["batch"]
        batch_counts[batch_num] = batch_counts.get(batch_num, 0) + 1

    return {
        "total_processed": len(results),
        "expected_count": expected_count,
        "all_processed": len(results) == expected_count,
        "batches": batch_counts,
    }


def format_report(summary: dict) -> str:
    """Format a summary as a human-readable string."""
    lines = [
        f"Total processed: {summary['total_processed']} / {summary['expected_count']}",
        f"Complete: {'Yes' if summary['all_processed'] else 'No'}",
        "Batch breakdown:",
    ]
    for batch_num, count in sorted(summary["batches"].items()):
        lines.append(f"  Batch {batch_num}: {count} records")
    return "\n".join(lines)
