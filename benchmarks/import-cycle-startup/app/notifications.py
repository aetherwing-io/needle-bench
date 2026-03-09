"""
Notification delivery module.

Handles sending emails and in-app notifications to users.
"""

import users
from config import NOTIFICATION_SENDER


# Default role for notification routing
_default_role = users.DEFAULT_ROLE


class NotificationService:
    """Sends notifications to users."""

    def __init__(self):
        self.pending = []
        self.sent_log = []

    def send_welcome_email(self, user_id: str, email: str):
        """Send a welcome email to a newly registered user."""
        message = {
            "type": "email",
            "from": NOTIFICATION_SENDER,
            "to": email,
            "subject": "Welcome to UserService!",
            "body": f"Hello! Your account (ID: {user_id}) has been created. "
                    f"Your default role is: {_default_role}",
        }
        self._deliver(message)

    def send_password_reset(self, user_id: str, email: str):
        """Send a password reset notification."""
        message = {
            "type": "email",
            "from": NOTIFICATION_SENDER,
            "to": email,
            "subject": "Password Reset Request",
            "body": f"A password reset was requested for account {user_id}.",
        }
        self._deliver(message)

    def send_login_alert(self, user_id: str, email: str, ip_address: str):
        """Send a login alert notification."""
        message = {
            "type": "email",
            "from": NOTIFICATION_SENDER,
            "to": email,
            "subject": "New Login Detected",
            "body": f"Account {user_id} logged in from {ip_address}.",
        }
        self._deliver(message)

    def notify_account_locked(self, user_id: str, email: str):
        """Notify user that their account has been locked."""
        message = {
            "type": "email",
            "from": NOTIFICATION_SENDER,
            "to": email,
            "subject": "Account Locked",
            "body": f"Your account {user_id} has been locked due to "
                    f"too many failed login attempts.",
        }
        self._deliver(message)

    def _deliver(self, message: dict):
        """Simulate message delivery (in production, this would send real emails)."""
        self.sent_log.append(message)
        print(f"[NOTIFY] {message['type']}: {message['subject']} -> {message['to']}")

    def get_sent_count(self) -> int:
        """Return the number of notifications sent."""
        return len(self.sent_log)
