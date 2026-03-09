#!/bin/sh
# Test for type-coercion-comparison bug
# Verifies that filtering by minimum rating returns the correct products.

set -e

cd /app

# Start the server in background
node server.js &
SERVER_PID=$!
sleep 2

cleanup() {
    kill $SERVER_PID 2>/dev/null || true
}
trap cleanup EXIT

FAIL=0

# Test 1: Filter by min_rating=4 should return products with rating >= 4
count=$(curl -s "http://localhost:3000/products?min_rating=4" | node -e "
const chunks = [];
process.stdin.on('data', c => chunks.push(c));
process.stdin.on('end', () => {
  const data = JSON.parse(Buffer.concat(chunks).toString());
  console.log(data.count);
});
")

# Products with rating >= 4: ids 1(4), 2(5), 4(4), 6(5), 9(5), 10(4), 13(4), 14(4), 15(5) = 9
if [ "$count" != "9" ]; then
    echo "FAIL: min_rating=4 should return 9 products, got $count"
    FAIL=1
fi

# Test 2: Filter by min_rating=5 should return only 5-star products
count2=$(curl -s "http://localhost:3000/products?min_rating=5" | node -e "
const chunks = [];
process.stdin.on('data', c => chunks.push(c));
process.stdin.on('end', () => {
  const data = JSON.parse(Buffer.concat(chunks).toString());
  console.log(data.count);
});
")

# Products with rating >= 5: ids 2(5), 6(5), 9(5), 15(5) = 4
if [ "$count2" != "4" ]; then
    echo "FAIL: min_rating=5 should return 4 products, got $count2"
    FAIL=1
fi

# Test 3: Filter by min_rating=3 should return products with rating >= 3
count3=$(curl -s "http://localhost:3000/products?min_rating=3" | node -e "
const chunks = [];
process.stdin.on('data', c => chunks.push(c));
process.stdin.on('end', () => {
  const data = JSON.parse(Buffer.concat(chunks).toString());
  console.log(data.count);
});
")

# Products with rating >= 3: ids 1(4),2(5),3(3),4(4),6(5),7(3),9(5),10(4),11(3),13(4),14(4),15(5) = 12
if [ "$count3" != "12" ]; then
    echo "FAIL: min_rating=3 should return 12 products, got $count3"
    FAIL=1
fi

if [ $FAIL -eq 0 ]; then
    echo "PASS: Rating filter returns correct results"
    exit 0
else
    exit 1
fi
