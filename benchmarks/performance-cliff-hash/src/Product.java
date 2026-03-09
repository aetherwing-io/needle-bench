/**
 * Represents a product in the inventory system.
 * Products are keyed by SKU (Stock Keeping Unit).
 */
public class Product {
    private final String sku;
    private final String name;
    private final String category;
    private final double price;
    private final int quantity;

    public Product(String sku, String name, String category, double price, int quantity) {
        this.sku = sku;
        this.name = name;
        this.category = category;
        this.price = price;
        this.quantity = quantity;
    }

    public String getSku() { return sku; }
    public String getName() { return name; }
    public String getCategory() { return category; }
    public double getPrice() { return price; }
    public int getQuantity() { return quantity; }

    /**
     * Custom hashCode for use in our ProductCache.
     * Groups products by their category for cache-friendly bucket placement.
     */
    @Override
    public int hashCode() {
        return category.hashCode();
    }

    @Override
    public boolean equals(Object obj) {
        if (this == obj) return true;
        if (obj == null || getClass() != obj.getClass()) return false;
        Product other = (Product) obj;
        return sku.equals(other.sku);
    }

    @Override
    public String toString() {
        return String.format("Product{sku='%s', name='%s', category='%s', price=%.2f, qty=%d}",
                sku, name, category, price, quantity);
    }
}
