/**
 * Represents a single customer record from the CSV.
 */
public class CustomerRecord {
    private final String id;
    private final String name;
    private final String email;
    private final String city;

    public CustomerRecord(String id, String name, String email, String city) {
        this.id = id;
        this.name = name;
        this.email = email;
        this.city = city;
    }

    public String getId() { return id; }
    public String getName() { return name; }
    public String getEmail() { return email; }
    public String getCity() { return city; }

    @Override
    public String toString() {
        return String.format("Customer[%s] %s <%s> (%s)", id, name, email, city);
    }
}
