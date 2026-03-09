import java.io.*;
import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.Map;
import java.util.TreeMap;

/**
 * Generates a customer summary report from parsed records.
 */
public class ReportGenerator {

    /**
     * Write a summary report to the given output path.
     * The report lists each customer and a city breakdown.
     */
    public void generateReport(List<CustomerRecord> customers, String outputPath)
            throws IOException {

        Map<String, Integer> cityCounts = new TreeMap<>();
        for (CustomerRecord c : customers) {
            cityCounts.merge(c.getCity(), 1, Integer::sum);
        }

        try (PrintWriter writer = new PrintWriter(
                new OutputStreamWriter(new FileOutputStream(outputPath), StandardCharsets.UTF_8))) {

            writer.println("=== Customer Summary Report ===");
            writer.println();
            writer.printf("Total customers: %d%n", customers.size());
            writer.println();

            writer.println("--- Customer List ---");
            for (CustomerRecord c : customers) {
                writer.printf("  [%s] %s <%s> - %s%n",
                    c.getId(), c.getName(), c.getEmail(), c.getCity());
            }
            writer.println();

            writer.println("--- City Breakdown ---");
            for (Map.Entry<String, Integer> entry : cityCounts.entrySet()) {
                writer.printf("  %s: %d customer(s)%n", entry.getKey(), entry.getValue());
            }
            writer.println();
            writer.println("=== End of Report ===");
        }
    }
}
