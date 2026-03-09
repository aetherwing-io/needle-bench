public class Main {
    public static void main(String[] args) {
        if (args.length < 1) {
            System.err.println("Usage: perftest <command>");
            System.err.println();
            System.err.println("Commands:");
            System.err.println("  bench          Run performance benchmark");
            System.err.println("  bench-small    Run benchmark with small dataset");
            System.err.println("  bench-large    Run benchmark with large dataset");
            System.err.println("  analyze        Analyze hash distribution");
            System.exit(1);
        }

        switch (args[0]) {
            case "bench":
                runBenchmark(10000);
                break;
            case "bench-small":
                runBenchmark(100);
                break;
            case "bench-large":
                runBenchmark(50000);
                break;
            case "analyze":
                analyzeDistribution();
                break;
            default:
                System.err.println("Unknown command: " + args[0]);
                System.exit(1);
        }
    }

    private static void runBenchmark(int size) {
        System.out.printf("=== Performance Benchmark (n=%d) ===%n", size);

        ProductCache cache = new ProductCache();

        // Insert products
        long insertStart = System.nanoTime();
        for (int i = 0; i < size; i++) {
            Product p = ProductGenerator.generate(i);
            cache.put(p);
        }
        long insertEnd = System.nanoTime();
        double insertMs = (insertEnd - insertStart) / 1_000_000.0;
        System.out.printf("Insert %d items: %.2f ms%n", size, insertMs);

        // Lookup all products
        long lookupStart = System.nanoTime();
        int found = 0;
        for (int i = 0; i < size; i++) {
            String sku = ProductGenerator.skuForIndex(i);
            Product p = cache.get(sku);
            if (p != null) found++;
        }
        long lookupEnd = System.nanoTime();
        double lookupMs = (lookupEnd - lookupStart) / 1_000_000.0;
        System.out.printf("Lookup %d items: %.2f ms (found %d)%n", size, lookupMs, found);

        // Performance threshold: at 10k entries, lookups should complete in < 500ms
        // With O(1) hash map, this is trivial. With O(n) degradation, it takes seconds.
        double avgLookupUs = (lookupMs * 1000) / size;
        System.out.printf("Average lookup: %.2f us%n", avgLookupUs);

        if (size >= 10000) {
            // At 10k entries, avg lookup should be < 50us with proper hashing
            // With hash collisions causing O(n), avg will be > 200us
            if (avgLookupUs > 100) {
                System.out.printf("FAIL: Average lookup %.2f us exceeds 100us threshold%n", avgLookupUs);
                System.out.printf("      This indicates O(n) degradation from hash collisions%n");
                System.exit(1);
            } else {
                System.out.printf("OK: Average lookup %.2f us is within acceptable range%n", avgLookupUs);
            }
        }
    }

    private static void analyzeDistribution() {
        System.out.println("=== Hash Distribution Analysis ===");

        ProductCache cache = new ProductCache();
        int size = 10000;

        for (int i = 0; i < size; i++) {
            Product p = ProductGenerator.generate(i);
            cache.put(p);
        }

        cache.printDistribution();
    }
}
