/**
 * Result of a transfer operation.
 */
public class TransferResult {
    private final boolean success;
    private final String message;

    public TransferResult(boolean success, String message) {
        this.success = success;
        this.message = message;
    }

    public boolean isSuccess() {
        return success;
    }

    public String getMessage() {
        return message;
    }

    @Override
    public String toString() {
        return String.format("TransferResult[success=%s, message=%s]", success, message);
    }
}
