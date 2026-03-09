# encoding-mojibake

A Java command-line tool that reads a CSV file of customer records and generates a summary report. The CSV contains international customer names with accented and non-ASCII characters (e.g., "Renee" with accents, umlauts, CJK characters). When the tool processes the file, non-ASCII characters in customer names are corrupted in the output -- accented letters become garbled multi-byte sequences. ASCII-only names display correctly.

## Symptoms

- Running `test.sh` shows that customer names with non-ASCII characters are mangled
- Characters like e-acute, u-umlaut, and CJK ideographs are replaced with wrong characters
- The report file contains mojibake (garbled text) for international names
- ASCII-only customer names like "John Smith" are unaffected

## Bug description

The application reads a UTF-8 encoded CSV file but decodes it incorrectly, causing multi-byte UTF-8 sequences to be misinterpreted. The corruption is consistent and deterministic -- the same wrong characters appear every time.

## Difficulty

Easy

## Expected turns

3-5
