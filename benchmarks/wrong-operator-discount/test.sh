#!/bin/sh
# Test for wrong-operator-discount bug
# Verifies that percentage discounts are calculated correctly.

set -e

cd /app

FAIL=0

# Test 1: 20% discount on known cart
result=$(python -c "
from cart import ShoppingCart
from catalog import get_product_catalog
from pricing import calculate_totals

catalog = get_product_catalog()
cart = ShoppingCart()
cart.add_item(catalog['LAPTOP'], 1)   # 999.99
cart.add_item(catalog['MOUSE'], 2)    # 29.99 * 2 = 59.98
cart.add_item(catalog['CABLE'], 3)    # 12.99 * 3 = 38.97
# subtotal = 1098.94

cart.apply_discount('SAVE20')  # 20% off
totals = calculate_totals(cart)
print(f\"{totals['discount_amount']:.2f}\")
")

# 20% of 1098.94 = 219.79
if [ "$result" != "219.79" ]; then
    echo "FAIL: 20% discount on \$1098.94 should be \$219.79, got \$$result"
    FAIL=1
fi

# Test 2: Verify final total
total=$(python -c "
from cart import ShoppingCart
from catalog import get_product_catalog
from pricing import calculate_totals

catalog = get_product_catalog()
cart = ShoppingCart()
cart.add_item(catalog['LAPTOP'], 1)
cart.add_item(catalog['MOUSE'], 2)
cart.add_item(catalog['CABLE'], 3)

cart.apply_discount('SAVE20')
totals = calculate_totals(cart)
print(f\"{totals['total']:.2f}\")
")

# 1098.94 - 219.79 = 879.15
if [ "$total" != "879.15" ]; then
    echo "FAIL: Total after 20% discount should be \$879.15, got \$$total"
    FAIL=1
fi

# Test 3: 10% discount on single item
result3=$(python -c "
from cart import ShoppingCart
from catalog import get_product_catalog
from pricing import calculate_totals

catalog = get_product_catalog()
cart = ShoppingCart()
cart.add_item(catalog['MONITOR'], 1)  # 449.99

cart.apply_discount('SAVE10')
totals = calculate_totals(cart)
print(f\"{totals['discount_amount']:.2f}\")
")

# 10% of 449.99 = 45.00
if [ "$result3" != "45.00" ]; then
    echo "FAIL: 10% discount on \$449.99 should be \$45.00, got \$$result3"
    FAIL=1
fi

# Test 4: No discount — total equals subtotal
result4=$(python -c "
from cart import ShoppingCart
from catalog import get_product_catalog
from pricing import calculate_totals

catalog = get_product_catalog()
cart = ShoppingCart()
cart.add_item(catalog['KEYBOARD'], 1)  # 89.99

totals = calculate_totals(cart)
print(f\"{totals['total']:.2f}\")
")

if [ "$result4" != "89.99" ]; then
    echo "FAIL: Total with no discount should be \$89.99, got \$$result4"
    FAIL=1
fi

if [ $FAIL -eq 0 ]; then
    echo "PASS: Discount calculations are correct"
    exit 0
else
    exit 1
fi
