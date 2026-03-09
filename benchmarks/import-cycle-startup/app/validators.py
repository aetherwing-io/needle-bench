"""
Input validation utilities.
"""

import re
from config import PASSWORD_MIN_LENGTH


def validate_email(email: str) -> bool:
    """Check if an email address is valid."""
    pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
    return bool(re.match(pattern, email))


def validate_password(password: str) -> bool:
    """Check if a password meets minimum requirements."""
    if len(password) < PASSWORD_MIN_LENGTH:
        return False
    has_upper = any(c.isupper() for c in password)
    has_lower = any(c.islower() for c in password)
    has_digit = any(c.isdigit() for c in password)
    return has_upper and has_lower and has_digit


def validate_username(username: str) -> bool:
    """Check if a username is valid."""
    if not username or len(username) < 3 or len(username) > 30:
        return False
    return bool(re.match(r'^[a-zA-Z0-9_]+$', username))
