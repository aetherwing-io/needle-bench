import java.util.List;

/**
 * Customer Report Tool
 *
 * Reads customer records from a CSV file and generates a formatted
 * summary report. Supports international customer names.
 *
 * Usage: java Main <input.csv> <output.txt>
 */
public class Main {

    public static void main(String[] args) {
        if (args.length < 2) {
            System.err.println("Usage: java Main <input.csv> <output.txt>");
            System.exit(1);
        }

        String inputPath = args[0];
        String outputPath = args[1];

        try {
            System.out.println("Reading customers from: " + inputPath);
            CsvReader reader = new CsvReader();
            List<CustomerRecord> customers = reader.readCustomers(inputPath);
            System.out.printf("Loaded %d customer records.%n", customers.size());

            System.out.println("Generating report: " + outputPath);
            ReportGenerator generator = new ReportGenerator();
            generator.generateReport(customers, outputPath);

            System.out.println("Report generated successfully.");
        } catch (Exception e) {
            System.err.println("ERROR: " + e.getMessage());
            e.printStackTrace();
            System.exit(2);
        }
    }
}
