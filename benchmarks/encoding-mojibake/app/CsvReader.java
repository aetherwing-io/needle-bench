import java.io.*;
import java.nio.charset.Charset;
import java.util.ArrayList;
import java.util.List;

/**
 * Reads customer records from a CSV file.
 *
 * CSV format: id,name,email,city
 * First line is a header row that gets skipped.
 */
public class CsvReader {

    private final Charset charset;

    public CsvReader() {
        // Use ISO-8859-1 for broad compatibility with legacy systems
        this.charset = Charset.forName("ISO-8859-1");
    }

    /**
     * Parse all customer records from a CSV file.
     */
    public List<CustomerRecord> readCustomers(String filePath) throws IOException {
        List<CustomerRecord> records = new ArrayList<>();

        try (BufferedReader reader = new BufferedReader(
                new InputStreamReader(new FileInputStream(filePath), charset))) {

            String line;
            boolean headerSkipped = false;

            while ((line = reader.readLine()) != null) {
                if (!headerSkipped) {
                    headerSkipped = true;
                    continue;
                }

                line = line.trim();
                if (line.isEmpty()) continue;

                String[] fields = parseCsvLine(line);
                if (fields.length < 4) {
                    System.err.println("WARN: Skipping malformed line: " + line);
                    continue;
                }

                records.add(new CustomerRecord(
                    fields[0].trim(),
                    fields[1].trim(),
                    fields[2].trim(),
                    fields[3].trim()
                ));
            }
        }

        return records;
    }

    /**
     * Simple CSV line parser (handles quoted fields).
     */
    private String[] parseCsvLine(String line) {
        List<String> fields = new ArrayList<>();
        StringBuilder current = new StringBuilder();
        boolean inQuotes = false;

        for (int i = 0; i < line.length(); i++) {
            char c = line.charAt(i);

            if (c == '"') {
                inQuotes = !inQuotes;
            } else if (c == ',' && !inQuotes) {
                fields.add(current.toString());
                current = new StringBuilder();
            } else {
                current.append(c);
            }
        }
        fields.add(current.toString());

        return fields.toArray(new String[0]);
    }
}
