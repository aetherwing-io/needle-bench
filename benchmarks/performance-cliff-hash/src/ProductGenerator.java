/**
 * Generates deterministic product data for benchmarking.
 */
public class ProductGenerator {

    private static final String[] CATEGORIES = {
        "Electronics", "Clothing", "Books", "Home", "Garden",
        "Sports", "Toys", "Food", "Beauty", "Automotive"
    };

    private static final String[] ADJECTIVES = {
        "Premium", "Basic", "Deluxe", "Pro", "Ultra",
        "Classic", "Modern", "Eco", "Smart", "Compact"
    };

    private static final String[] NOUNS = {
        "Widget", "Gadget", "Device", "Tool", "Kit",
        "Pack", "Set", "Bundle", "Unit", "Module"
    };

    /**
     * Generate a deterministic SKU for a given index.
     */
    public static String skuForIndex(int index) {
        return String.format("SKU-%06d", index);
    }

    /**
     * Generate a product with deterministic fields based on index.
     */
    public static Product generate(int index) {
        String sku = skuForIndex(index);
        String category = CATEGORIES[index % CATEGORIES.length];
        String adjective = ADJECTIVES[(index / CATEGORIES.length) % ADJECTIVES.length];
        String noun = NOUNS[(index / (CATEGORIES.length * ADJECTIVES.length)) % NOUNS.length];
        String name = adjective + " " + noun;
        double price = 9.99 + (index % 100) * 5.0;
        int quantity = 1 + (index % 500);

        return new Product(sku, name, category, price, quantity);
    }
}
