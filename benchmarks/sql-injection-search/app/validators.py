"""
Input validation utilities.
"""
import re


def validate_product_id(product_id):
    """Validate that product_id is a positive integer."""
    try:
        pid = int(product_id)
        return pid > 0
    except (ValueError, TypeError):
        return False


def validate_price(price):
    """Validate that price is a positive number."""
    try:
        p = float(price)
        return p >= 0
    except (ValueError, TypeError):
        return False


def sanitize_string(s, max_length=200):
    """Basic string sanitization — truncate and strip."""
    if not isinstance(s, str):
        return ''
    return s.strip()[:max_length]


def validate_category(category):
    """Validate category is one of the known values."""
    valid = {'electronics', 'furniture', 'office', 'clothing', 'food'}
    return category in valid if category else True
