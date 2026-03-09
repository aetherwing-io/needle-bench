"""
UserService — main entry point.

A simple user management service with notification support.
"""

import sys
from users import UserManager


def main():
    print("Starting UserService...")

    manager = UserManager()

    # Register some test users
    print("\nRegistering users...")
    try:
        alice = manager.register("alice", "alice@example.com", "SecurePass1")
        print(f"  Registered: {alice}")

        bob = manager.register("bob_smith", "bob@example.com", "MyPassword2")
        print(f"  Registered: {bob}")
    except Exception as e:
        print(f"  ERROR during registration: {e}")
        return 1

    # List all users
    print("\nAll users:")
    for user in manager.list_users():
        print(f"  {user}")

    # Test authentication
    print("\nTesting authentication...")
    result = manager.authenticate("alice", "SecurePass1")
    if result:
        print(f"  alice authenticated: {result}")
    else:
        print("  alice authentication FAILED")
        return 1

    result = manager.authenticate("alice", "WrongPassword")
    if result is None:
        print("  Wrong password correctly rejected")
    else:
        print("  ERROR: Wrong password was accepted")
        return 1

    print("\nUserService started successfully.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
