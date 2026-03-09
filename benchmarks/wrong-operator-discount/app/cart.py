"""Shopping cart module."""

from typing import Optional


class CartItem:
    """An item in the shopping cart."""

    def __init__(self, product: dict, quantity: int):
        self.product = product
        self.quantity = quantity

    @property
    def subtotal(self) -> float:
        return round(self.product["price"] * self.quantity, 2)

    def to_dict(self) -> dict:
        return {
            "sku": self.product["sku"],
            "name": self.product["name"],
            "price": self.product["price"],
            "quantity": self.quantity,
            "subtotal": self.subtotal,
        }


class ShoppingCart:
    """Shopping cart with discount support."""

    def __init__(self):
        self.items: list[CartItem] = []
        self.discount_code: Optional[str] = None
        self.discount_percent: float = 0.0

    def add_item(self, product: dict, quantity: int = 1) -> None:
        """Add a product to the cart."""
        for item in self.items:
            if item.product["sku"] == product["sku"]:
                item.quantity += quantity
                return
        self.items.append(CartItem(product, quantity))

    def remove_item(self, sku: str) -> bool:
        """Remove an item from the cart by SKU."""
        for i, item in enumerate(self.items):
            if item.product["sku"] == sku:
                self.items.pop(i)
                return True
        return False

    def apply_discount(self, code: str) -> bool:
        """Apply a discount code to the cart."""
        discounts = {
            "SAVE10": 10,
            "SAVE20": 20,
            "SAVE50": 50,
            "WELCOME5": 5,
        }
        if code in discounts:
            self.discount_code = code
            self.discount_percent = discounts[code]
            return True
        return False

    @property
    def subtotal(self) -> float:
        """Total before discounts."""
        return round(sum(item.subtotal for item in self.items), 2)

    def to_dict(self) -> dict:
        return {
            "items": [item.to_dict() for item in self.items],
            "item_count": len(self.items),
            "subtotal": self.subtotal,
            "discount_code": self.discount_code,
            "discount_percent": self.discount_percent,
        }
