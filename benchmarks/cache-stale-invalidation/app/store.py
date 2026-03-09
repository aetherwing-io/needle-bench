"""
In-memory product store (simulates a database).
"""
import uuid
import time


class ProductStore:
    def __init__(self):
        self._products = {}
        self._seed_data()

    def _seed_data(self):
        """Pre-populate with sample products."""
        samples = [
            {'name': 'Widget A', 'price': 9.99, 'category': 'widgets'},
            {'name': 'Widget B', 'price': 14.99, 'category': 'widgets'},
            {'name': 'Gadget X', 'price': 29.99, 'category': 'gadgets'},
            {'name': 'Gadget Y', 'price': 49.99, 'category': 'gadgets'},
            {'name': 'Tool Z', 'price': 19.99, 'category': 'tools'},
        ]
        for data in samples:
            self.create_product(data)

    def list_products(self):
        """Return all products as a list."""
        return list(self._products.values())

    def get_product(self, product_id):
        """Get a single product by ID."""
        return self._products.get(product_id)

    def create_product(self, data):
        """Create a new product and return it."""
        product_id = str(uuid.uuid4())[:8]
        product = {
            'id': product_id,
            'name': data.get('name', 'Unnamed'),
            'price': float(data.get('price', 0)),
            'category': data.get('category', 'uncategorized'),
            'created_at': time.time(),
            'updated_at': time.time(),
        }
        self._products[product_id] = product
        return product

    def update_product(self, product_id, data):
        """Update an existing product. Returns updated product."""
        product = self._products.get(product_id)
        if product is None:
            return None

        if 'name' in data:
            product['name'] = data['name']
        if 'price' in data:
            product['price'] = float(data['price'])
        if 'category' in data:
            product['category'] = data['category']

        product['updated_at'] = time.time()
        return product

    def delete_product(self, product_id):
        """Delete a product by ID."""
        if product_id in self._products:
            del self._products[product_id]
            return True
        return False
