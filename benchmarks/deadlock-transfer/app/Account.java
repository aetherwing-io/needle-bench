/**
 * Represents a bank account with thread-safe balance operations.
 */
public class Account {
    private final String id;
    private final String owner;
    private double balance;
    private final Object lock = new Object();

    public Account(String id, String owner, double initialBalance) {
        this.id = id;
        this.owner = owner;
        this.balance = initialBalance;
    }

    public String getId() {
        return id;
    }

    public String getOwner() {
        return owner;
    }

    public double getBalance() {
        synchronized (lock) {
            return balance;
        }
    }

    public Object getLock() {
        return lock;
    }

    /**
     * Withdraw amount from account. Must be called while holding lock.
     */
    public boolean withdraw(double amount) {
        if (amount <= 0) throw new IllegalArgumentException("Amount must be positive");
        if (balance < amount) return false;
        balance -= amount;
        return true;
    }

    /**
     * Deposit amount into account. Must be called while holding lock.
     */
    public void deposit(double amount) {
        if (amount <= 0) throw new IllegalArgumentException("Amount must be positive");
        balance += amount;
    }

    @Override
    public String toString() {
        return String.format("Account[%s, owner=%s, balance=%.2f]", id, owner, balance);
    }
}
