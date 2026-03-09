import java.util.LinkedList;

/**
 * A hash map-based cache for Product objects, keyed by SKU.
 * Uses separate chaining for collision resolution.
 *
 * This is a simplified HashMap implementation to demonstrate
 * the importance of good hash functions. The cache uses the
 * Product's hashCode() to determine bucket placement.
 */
public class ProductCache {

    private static final int DEFAULT_CAPACITY = 1024;
    private static final float LOAD_FACTOR = 0.75f;

    private LinkedList<Entry>[] buckets;
    private int size;
    private int capacity;

    private static class Entry {
        String key;
        Product value;

        Entry(String key, Product value) {
            this.key = key;
            this.value = value;
        }
    }

    @SuppressWarnings("unchecked")
    public ProductCache() {
        this.capacity = DEFAULT_CAPACITY;
        this.buckets = new LinkedList[capacity];
        this.size = 0;
    }

    /**
     * Store a product in the cache.
     */
    public void put(Product product) {
        int index = getBucketIndex(product);
        if (buckets[index] == null) {
            buckets[index] = new LinkedList<>();
        }

        // Update existing entry if SKU matches
        for (Entry entry : buckets[index]) {
            if (entry.key.equals(product.getSku())) {
                entry.value = product;
                return;
            }
        }

        buckets[index].add(new Entry(product.getSku(), product));
        size++;

        // Resize if needed
        if ((float) size / capacity > LOAD_FACTOR) {
            resize();
        }
    }

    /**
     * Retrieve a product by SKU.
     */
    public Product get(String sku) {
        // To look up by SKU, we need to compute the bucket index.
        // We create a dummy product to use the same hashing path.
        // This mirrors the put() path which uses Product.hashCode().
        int index = getBucketIndexForSku(sku);
        if (buckets[index] == null) {
            return null;
        }

        for (Entry entry : buckets[index]) {
            if (entry.key.equals(sku)) {
                return entry.value;
            }
        }
        return null;
    }

    /**
     * Get the bucket index for a product using its hashCode.
     */
    private int getBucketIndex(Product product) {
        int hash = product.hashCode();
        return Math.abs(hash) % capacity;
    }

    /**
     * Get the bucket index for a SKU lookup.
     * Must match the same bucket as put() for the same SKU.
     * Reconstructs the product to compute the correct hash.
     */
    private int getBucketIndexForSku(String sku) {
        // Extract the index from SKU format "SKU-NNNNNN"
        int index;
        try {
            index = Integer.parseInt(sku.substring(4));
        } catch (NumberFormatException e) {
            return 0;
        }

        // Reconstruct the product to compute its hashCode
        Product dummy = ProductGenerator.generate(index);
        return getBucketIndex(dummy);
    }

    /**
     * Get the number of entries in the cache.
     */
    public int size() {
        return size;
    }

    /**
     * Print bucket distribution statistics.
     */
    public void printDistribution() {
        int nonEmpty = 0;
        int maxChain = 0;
        long totalChain = 0;

        for (int i = 0; i < capacity; i++) {
            if (buckets[i] != null && !buckets[i].isEmpty()) {
                nonEmpty++;
                int chainLen = buckets[i].size();
                totalChain += chainLen;
                if (chainLen > maxChain) {
                    maxChain = chainLen;
                }
            }
        }

        System.out.printf("Capacity: %d%n", capacity);
        System.out.printf("Size: %d%n", size);
        System.out.printf("Non-empty buckets: %d (%.1f%%)%n", nonEmpty, (nonEmpty * 100.0 / capacity));
        System.out.printf("Max chain length: %d%n", maxChain);
        if (nonEmpty > 0) {
            System.out.printf("Avg chain length: %.1f%n", (double) totalChain / nonEmpty);
        }

        if (maxChain > 100) {
            System.out.printf("WARNING: max chain length %d indicates severe hash collisions%n", maxChain);
        }
    }

    @SuppressWarnings("unchecked")
    private void resize() {
        int newCapacity = capacity * 2;
        LinkedList<Entry>[] newBuckets = new LinkedList[newCapacity];

        for (int i = 0; i < capacity; i++) {
            if (buckets[i] != null) {
                for (Entry entry : buckets[i]) {
                    int newIndex = Math.abs(entry.value.hashCode()) % newCapacity;
                    if (newBuckets[newIndex] == null) {
                        newBuckets[newIndex] = new LinkedList<>();
                    }
                    newBuckets[newIndex].add(entry);
                }
            }
        }

        buckets = newBuckets;
        capacity = newCapacity;
    }
}
