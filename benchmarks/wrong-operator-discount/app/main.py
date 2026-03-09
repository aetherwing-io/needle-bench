"""Shopping cart CLI with discount support."""

import sys
import json
from cart import ShoppingCart
from catalog import get_product_catalog
from pricing import calculate_totals


def main():
    catalog = get_product_catalog()
    cart = ShoppingCart()

    # Simulate a typical order
    cart.add_item(catalog["LAPTOP"], 1)
    cart.add_item(catalog["MOUSE"], 2)
    cart.add_item(catalog["CABLE"], 3)

    # Apply a 20% discount code
    discount_code = sys.argv[1] if len(sys.argv) > 1 else "SAVE20"
    cart.apply_discount(discount_code)

    totals = calculate_totals(cart)
    print(json.dumps(totals, indent=2))

    return totals


if __name__ == "__main__":
    main()
