# wrong-operator-discount

## Project

A Python shopping cart with discount code support. Calculates subtotals, applies percentage-based discounts, and produces order totals.

## Symptoms

When applying a discount code (e.g., "SAVE20" for 20% off), the calculated discount amount is wildly incorrect. A 20% discount on a ~$1,099 order produces a discount of ~$1,099 instead of ~$220, resulting in a total of $0. Customers either get free orders or the total is nonsensically wrong.

## Bug description

The discount calculation formula produces incorrect results. The arithmetic is wrong in how the percentage is applied to the subtotal. The function compiles and runs without errors but the math is incorrect.

## Difficulty

Easy

## Expected turns

2-4
