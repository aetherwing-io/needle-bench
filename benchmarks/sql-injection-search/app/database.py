"""
SQLite database wrapper.
"""
import sqlite3
import os


class Database:
    def __init__(self, db_path):
        self.db_path = db_path
        self.conn = None

    def _connect(self):
        if self.conn is None:
            self.conn = sqlite3.connect(self.db_path)
            self.conn.row_factory = sqlite3.Row

    def initialize(self):
        """Create tables and seed data."""
        self._connect()
        cursor = self.conn.cursor()

        cursor.execute('''
            CREATE TABLE IF NOT EXISTS products (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name TEXT NOT NULL,
                description TEXT,
                price REAL NOT NULL,
                category TEXT NOT NULL,
                in_stock INTEGER DEFAULT 1
            )
        ''')

        cursor.execute('''
            CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                username TEXT NOT NULL UNIQUE,
                password_hash TEXT NOT NULL,
                email TEXT NOT NULL,
                role TEXT DEFAULT 'user',
                api_key TEXT
            )
        ''')

        # Check if data already exists
        count = cursor.execute('SELECT COUNT(*) FROM products').fetchone()[0]
        if count == 0:
            self._seed_products(cursor)
            self._seed_users(cursor)

        self.conn.commit()

    def _seed_products(self, cursor):
        products = [
            ('Wireless Mouse', 'Ergonomic wireless mouse with USB receiver', 29.99, 'electronics', 1),
            ('Mechanical Keyboard', 'Cherry MX Blue switches, full size', 89.99, 'electronics', 1),
            ('USB-C Hub', '7-port USB-C hub with HDMI', 45.99, 'electronics', 1),
            ('Standing Desk', 'Electric adjustable standing desk 60x30', 399.99, 'furniture', 1),
            ('Monitor Arm', 'Single monitor arm, VESA compatible', 79.99, 'furniture', 1),
            ('Desk Lamp', 'LED desk lamp with USB charging port', 34.99, 'furniture', 0),
            ('Webcam HD', '1080p webcam with built-in microphone', 59.99, 'electronics', 1),
            ('Notebook', 'Leather-bound notebook, 200 pages', 12.99, 'office', 1),
            ('Pen Set', 'Professional gel pen set, 10 colors', 8.99, 'office', 1),
            ('Cable Organizer', 'Silicon cable management clips, 12 pack', 9.99, 'office', 1),
        ]
        cursor.executemany(
            'INSERT INTO products (name, description, price, category, in_stock) VALUES (?, ?, ?, ?, ?)',
            products
        )

    def _seed_users(self, cursor):
        users = [
            ('admin', 'pbkdf2:sha256:admin_hash_xxx', 'admin@example.com', 'admin', 'sk-admin-secret-key-123'),
            ('alice', 'pbkdf2:sha256:alice_hash_xxx', 'alice@example.com', 'user', 'sk-alice-key-456'),
            ('bob', 'pbkdf2:sha256:bob_hash_xxx', 'bob@example.com', 'user', None),
        ]
        cursor.executemany(
            'INSERT INTO users (username, password_hash, email, role, api_key) VALUES (?, ?, ?, ?, ?)',
            users
        )

    def execute(self, query, params=None):
        """Execute a parameterized query and return rows as dicts."""
        self._connect()
        cursor = self.conn.cursor()
        if params:
            cursor.execute(query, params)
        else:
            cursor.execute(query)
        columns = [col[0] for col in cursor.description] if cursor.description else []
        return [dict(zip(columns, row)) for row in cursor.fetchall()]

    def execute_raw(self, query):
        """Execute a raw SQL query string. Used for dynamic queries."""
        self._connect()
        cursor = self.conn.cursor()
        cursor.execute(query)
        columns = [col[0] for col in cursor.description] if cursor.description else []
        return [dict(zip(columns, row)) for row in cursor.fetchall()]

    def close(self):
        if self.conn:
            self.conn.close()
            self.conn = None
