"""Configuration for the batch processor."""

# Default batch processing settings
DEFAULT_BATCH_SIZE = 10
MAX_BATCH_SIZE = 1000
MIN_BATCH_SIZE = 1

# Processing options
UPPERCASE_VALUES = True
INCLUDE_METADATA = True


def validate_batch_size(size: int) -> int:
    """Validate and clamp batch size to allowed range."""
    if size < MIN_BATCH_SIZE:
        return MIN_BATCH_SIZE
    if size > MAX_BATCH_SIZE:
        return MAX_BATCH_SIZE
    return size
