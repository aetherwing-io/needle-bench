/**
 * Service that handles fund transfers between accounts.
 * Uses locking to ensure consistency during concurrent transfers.
 */
public class TransferService {
    private int completedTransfers = 0;
    private int failedTransfers = 0;

    /**
     * Transfer funds from one account to another.
     * BUG: Locks accounts in parameter order (from, then to).
     * When two concurrent transfers go in opposite directions
     * (A->B and B->A), they lock in opposite order -> deadlock.
     */
    public TransferResult transfer(Account from, Account to, double amount) {
        if (from.getId().equals(to.getId())) {
            return new TransferResult(false, "Cannot transfer to same account");
        }

        if (amount <= 0) {
            return new TransferResult(false, "Amount must be positive");
        }

        // Lock the source account first, then the destination
        synchronized (from.getLock()) {
            // Small delay to increase deadlock probability
            try { Thread.sleep(1); } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }

            synchronized (to.getLock()) {
                if (from.getBalance() < amount) {
                    failedTransfers++;
                    return new TransferResult(false,
                        String.format("Insufficient funds: %.2f available, %.2f requested",
                            from.getBalance(), amount));
                }

                from.withdraw(amount);
                to.deposit(amount);
                completedTransfers++;

                return new TransferResult(true,
                    String.format("Transferred %.2f from %s to %s",
                        amount, from.getId(), to.getId()));
            }
        }
    }

    public int getCompletedTransfers() {
        return completedTransfers;
    }

    public int getFailedTransfers() {
        return failedTransfers;
    }
}
