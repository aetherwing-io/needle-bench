# performance-cliff-hash

## Project

A Java inventory cache system (`perftest`) that stores products in a custom hash map implementation. Products are keyed by SKU (Stock Keeping Unit) and stored using separate-chaining hash buckets. The system includes a benchmark harness that measures insert and lookup performance at various dataset sizes, plus a distribution analyzer that shows how entries are spread across buckets.

## Symptoms

With small datasets (100 items), the cache performs well. At 10,000 entries, lookup performance degrades dramatically. The average lookup time exceeds the acceptable threshold by a large margin. The distribution analyzer shows that most entries are concentrated in a very small number of buckets, with extremely long chain lengths, while the vast majority of buckets remain empty.

## Bug description

The hash function used to distribute products across buckets does not use the unique identifier of each product. Instead, it hashes on a field that is shared among many products, causing massive collision rates. At small scales the collision chains are short enough to appear performant. At scale, the O(n) chain traversal dominates lookup time. The fix requires using the unique product key in the hash computation so entries distribute evenly across buckets.

## Difficulty

Hard

## Expected turns

10-15
