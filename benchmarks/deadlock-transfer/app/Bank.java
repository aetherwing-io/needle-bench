import java.util.Map;
import java.util.HashMap;
import java.util.Collection;

/**
 * Bank holds a set of accounts and provides lookup.
 */
public class Bank {
    private final Map<String, Account> accounts = new HashMap<>();
    private final TransferService transferService;

    public Bank() {
        this.transferService = new TransferService();
    }

    public Account createAccount(String id, String owner, double initialBalance) {
        Account account = new Account(id, owner, initialBalance);
        accounts.put(id, account);
        return account;
    }

    public Account getAccount(String id) {
        Account account = accounts.get(id);
        if (account == null) {
            throw new IllegalArgumentException("Account not found: " + id);
        }
        return account;
    }

    public TransferResult transfer(String fromId, String toId, double amount) {
        Account from = getAccount(fromId);
        Account to = getAccount(toId);
        return transferService.transfer(from, to, amount);
    }

    public TransferService getTransferService() {
        return transferService;
    }

    public Collection<Account> getAllAccounts() {
        return accounts.values();
    }

    public double getTotalBalance() {
        return accounts.values().stream()
            .mapToDouble(Account::getBalance)
            .sum();
    }
}
