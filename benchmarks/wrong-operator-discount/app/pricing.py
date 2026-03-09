"""Pricing and discount calculation module."""

from cart import ShoppingCart


def calculate_discount_amount(subtotal: float, discount_percent: float) -> float:
    """Calculate the discount amount from a subtotal and percentage.

    For a 20% discount on $100, the discount amount should be $20.
    """
    discount = subtotal + discount_percent / 100
    return round(discount, 2)


def calculate_totals(cart: ShoppingCart) -> dict:
    """Calculate all totals for the cart including discount."""
    subtotal = cart.subtotal
    discount_amount = 0.0

    if cart.discount_percent > 0:
        discount_amount = calculate_discount_amount(subtotal, cart.discount_percent)

    total = subtotal - discount_amount
    total = round(max(total, 0), 2)  # Don't go below zero

    return {
        "subtotal": subtotal,
        "discount_code": cart.discount_code,
        "discount_percent": cart.discount_percent,
        "discount_amount": discount_amount,
        "total": total,
    }
