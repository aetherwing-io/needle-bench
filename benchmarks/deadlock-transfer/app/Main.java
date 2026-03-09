import java.util.concurrent.*;
import java.util.ArrayList;
import java.util.List;

/**
 * Entry point — runs concurrent transfers and reports results.
 * Usage: java Main [numTransfers] [numThreads]
 */
public class Main {
    public static void main(String[] args) throws Exception {
        int numTransfers = args.length > 0 ? Integer.parseInt(args[0]) : 100;
        int numThreads = args.length > 1 ? Integer.parseInt(args[1]) : 4;
        int timeoutSeconds = args.length > 2 ? Integer.parseInt(args[2]) : 10;

        Bank bank = new Bank();
        bank.createAccount("A", "Alice", 10000.0);
        bank.createAccount("B", "Bob", 10000.0);
        bank.createAccount("C", "Carol", 10000.0);

        double initialTotal = bank.getTotalBalance();
        System.out.printf("Initial total balance: %.2f%n", initialTotal);
        System.out.printf("Running %d transfers across %d threads (timeout: %ds)%n",
            numTransfers, numThreads, timeoutSeconds);

        ExecutorService executor = Executors.newFixedThreadPool(numThreads);
        List<Future<?>> futures = new ArrayList<>();

        // Submit transfers — mix of directions to trigger deadlock
        String[][] pairs = {
            {"A", "B"}, {"B", "A"},
            {"B", "C"}, {"C", "B"},
            {"A", "C"}, {"C", "A"},
        };

        for (int i = 0; i < numTransfers; i++) {
            String[] pair = pairs[i % pairs.length];
            final String from = pair[0];
            final String to = pair[1];
            final double amount = 10.0 + (i % 50);

            futures.add(executor.submit(() -> {
                TransferResult result = bank.transfer(from, to, amount);
                if (!result.isSuccess()) {
                    System.err.println("Transfer failed: " + result.getMessage());
                }
            }));
        }

        executor.shutdown();
        boolean finished = executor.awaitTermination(timeoutSeconds, TimeUnit.SECONDS);

        if (!finished) {
            System.err.println("ERROR: Transfers did not complete within timeout — possible deadlock");
            executor.shutdownNow();
            System.exit(2);
        }

        double finalTotal = bank.getTotalBalance();
        System.out.printf("Final total balance: %.2f%n", finalTotal);
        System.out.printf("Completed: %d, Failed: %d%n",
            bank.getTransferService().getCompletedTransfers(),
            bank.getTransferService().getFailedTransfers());

        if (Math.abs(finalTotal - initialTotal) > 0.01) {
            System.err.printf("ERROR: Balance mismatch! Initial=%.2f, Final=%.2f%n",
                initialTotal, finalTotal);
            System.exit(3);
        }

        System.out.println("All transfers completed successfully");
        System.exit(0);
    }
}
