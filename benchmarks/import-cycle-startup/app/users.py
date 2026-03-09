"""
User management module.

Handles user registration, authentication, and account operations.
"""

import hashlib
from config import MAX_LOGIN_ATTEMPTS
from validators import validate_email, validate_password, validate_username
from notifications import NotificationService


# Default role assigned to new users
DEFAULT_ROLE = "member"


class UserManager:
    """Manages user accounts."""

    def __init__(self):
        self._users = {}
        self._next_id = 1000
        self._notifier = NotificationService()

    def register(self, username: str, email: str, password: str) -> dict:
        """Register a new user account."""
        if not validate_username(username):
            raise ValueError(f"Invalid username: {username}")
        if not validate_email(email):
            raise ValueError(f"Invalid email: {email}")
        if not validate_password(password):
            raise ValueError("Password does not meet requirements")

        # Check for duplicate username
        for user in self._users.values():
            if user["username"] == username:
                raise ValueError(f"Username already taken: {username}")

        user_id = f"usr-{self._next_id}"
        self._next_id += 1

        user = {
            "id": user_id,
            "username": username,
            "email": email,
            "password_hash": self._hash_password(password),
            "role": DEFAULT_ROLE,
            "login_attempts": 0,
            "locked": False,
        }
        self._users[user_id] = user

        # Send welcome notification
        self._notifier.send_welcome_email(user_id, email)

        return {"id": user_id, "username": username, "email": email}

    def authenticate(self, username: str, password: str) -> dict | None:
        """Authenticate a user by username and password."""
        user = self._find_by_username(username)
        if not user:
            return None

        if user["locked"]:
            return None

        if user["password_hash"] == self._hash_password(password):
            user["login_attempts"] = 0
            return {"id": user["id"], "username": user["username"]}
        else:
            user["login_attempts"] += 1
            if user["login_attempts"] >= MAX_LOGIN_ATTEMPTS:
                user["locked"] = True
                self._notifier.notify_account_locked(user["id"], user["email"])
            return None

    def get_user(self, user_id: str) -> dict | None:
        """Retrieve a user by ID."""
        return self._users.get(user_id)

    def list_users(self) -> list[dict]:
        """Return all users (without password hashes)."""
        return [
            {"id": u["id"], "username": u["username"], "email": u["email"]}
            for u in self._users.values()
        ]

    def _find_by_username(self, username: str) -> dict | None:
        """Find a user by username."""
        for user in self._users.values():
            if user["username"] == username:
                return user
        return None

    @staticmethod
    def _hash_password(password: str) -> str:
        """Hash a password (simplified for demo)."""
        return hashlib.sha256(password.encode()).hexdigest()
