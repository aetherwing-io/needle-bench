"""Product catalog."""


def get_product_catalog() -> dict:
    """Return the product catalog as a dict keyed by SKU."""
    products = [
        {"sku": "LAPTOP", "name": "Business Laptop", "price": 999.99, "category": "electronics"},
        {"sku": "MOUSE", "name": "Wireless Mouse", "price": 29.99, "category": "electronics"},
        {"sku": "CABLE", "name": "USB-C Cable", "price": 12.99, "category": "accessories"},
        {"sku": "MONITOR", "name": "4K Monitor", "price": 449.99, "category": "electronics"},
        {"sku": "KEYBOARD", "name": "Mechanical Keyboard", "price": 89.99, "category": "electronics"},
        {"sku": "HEADSET", "name": "Noise-Canceling Headset", "price": 199.99, "category": "electronics"},
    ]
    return {p["sku"]: p for p in products}


def get_product(sku: str) -> dict | None:
    """Look up a single product by SKU."""
    catalog = get_product_catalog()
    return catalog.get(sku)
