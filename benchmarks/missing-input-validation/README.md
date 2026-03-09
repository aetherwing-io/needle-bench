# missing-input-validation

## Project

A TypeScript HTTP API for inventory management. Supports listing products, viewing individual items, and adjusting stock quantities via POST requests.

## Symptoms

The API accepts inventory adjustments that result in negative stock quantities. For example, a product with 100 units can have -200 applied, resulting in -100 units in stock. This creates impossible inventory states that cascade into billing errors and fulfillment failures downstream.

## Bug description

The quantity adjustment endpoint validates that the input is a number but does not check whether the resulting inventory would become negative. Any numeric adjustment is blindly applied. The fix requires rejecting adjustments that would drive inventory below zero.

## Difficulty

Easy

## Expected turns

3-5
