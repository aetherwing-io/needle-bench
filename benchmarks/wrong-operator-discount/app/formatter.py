"""Output formatting utilities."""


def format_currency(amount: float) -> str:
    """Format a number as USD currency."""
    return f"${amount:,.2f}"


def format_receipt(totals: dict) -> str:
    """Format cart totals as a receipt string."""
    lines = [
        "=" * 40,
        "ORDER SUMMARY",
        "=" * 40,
        f"Subtotal:     {format_currency(totals['subtotal'])}",
    ]

    if totals.get("discount_code"):
        lines.append(
            f"Discount:     -{format_currency(totals['discount_amount'])} "
            f"({totals['discount_code']}, {totals['discount_percent']}% off)"
        )

    lines.extend([
        "-" * 40,
        f"Total:        {format_currency(totals['total'])}",
        "=" * 40,
    ])

    return "\n".join(lines)
