#!/bin/sh
# Test for stale cache after writes.
# Verifies that reads after updates return fresh data, not stale cached values.

set -e

FAIL=0
PORT=18081

# Start server
python3 /app/app/server.py $PORT &
SERVER_PID=$!
sleep 1

if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo "FAIL: Server failed to start"
    exit 1
fi

# Step 1: Read product list to populate cache
LIST1=$(curl -s http://localhost:$PORT/products)
COUNT1=$(echo "$LIST1" | jq 'length')
echo "Initial product count: $COUNT1"

# Step 2: Get a specific product to cache it
FIRST_ID=$(echo "$LIST1" | jq -r '.[0].id')
PRODUCT1=$(curl -s http://localhost:$PORT/products/$FIRST_ID)
ORIG_NAME=$(echo "$PRODUCT1" | jq -r '.name')
echo "Original product name: $ORIG_NAME"

# Step 3: Update the product's name
curl -s -X PUT http://localhost:$PORT/products/$FIRST_ID \
    -H "Content-Type: application/json" \
    -d '{"name":"UPDATED-NAME-XYZ"}' >/dev/null

# Step 4: Read the product again — should get the updated name
PRODUCT2=$(curl -s http://localhost:$PORT/products/$FIRST_ID)
UPDATED_NAME=$(echo "$PRODUCT2" | jq -r '.name')
echo "Name after update: $UPDATED_NAME"

if [ "$UPDATED_NAME" != "UPDATED-NAME-XYZ" ]; then
    echo "FAIL: Product name is '$UPDATED_NAME' after update, expected 'UPDATED-NAME-XYZ'"
    echo "      Cache is serving stale data after write"
    FAIL=1
else
    echo "PASS: Product name correctly updated to '$UPDATED_NAME'"
fi

# Step 5: Create a new product
curl -s -X POST http://localhost:$PORT/products \
    -H "Content-Type: application/json" \
    -d '{"name":"New Product","price":99.99,"category":"test"}' >/dev/null

# Step 6: List products — should include the new product
LIST2=$(curl -s http://localhost:$PORT/products)
COUNT2=$(echo "$LIST2" | jq 'length')
echo "Product count after create: $COUNT2"

EXPECTED_COUNT=$((COUNT1 + 1))
if [ "$COUNT2" -ne "$EXPECTED_COUNT" ]; then
    echo "FAIL: Product list has $COUNT2 items after create, expected $EXPECTED_COUNT"
    echo "      List cache is stale — new product not visible"
    FAIL=1
else
    echo "PASS: Product list correctly shows $COUNT2 items"
fi

# Cleanup
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

if [ $FAIL -eq 0 ]; then
    echo "PASS: Cache correctly invalidated after writes"
    exit 0
else
    exit 1
fi
