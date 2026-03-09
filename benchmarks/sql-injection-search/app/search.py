"""
Search engine for product queries.
"""


class SearchEngine:
    def __init__(self, db):
        self.db = db

    def search(self, query, category=None):
        """
        Search for products matching the query string.
        Optionally filter by category.

        BUG: User input is concatenated directly into the SQL query string
        instead of using parameterized queries. This allows SQL injection.
        """
        sql = f"SELECT id, name, description, price, category, in_stock FROM products WHERE name LIKE '%{query}%'"

        if category:
            sql += f" AND category = '{category}'"

        sql += " ORDER BY name"

        return self.db.execute_raw(sql)

    def search_by_price_range(self, min_price, max_price):
        """Search products within a price range (safely parameterized)."""
        return self.db.execute(
            'SELECT * FROM products WHERE price BETWEEN ? AND ? ORDER BY price',
            (min_price, max_price)
        )

    def get_categories(self):
        """Get all distinct categories."""
        return self.db.execute('SELECT DISTINCT category FROM products ORDER BY category')
