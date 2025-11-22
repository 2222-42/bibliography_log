# Bibliography Log

![Go Tests](https://github.com/2222-42/bibliography_log/actions/workflows/test.yml/badge.svg)
![Lint](https://github.com/2222-42/bibliography_log/actions/workflows/lint.yml/badge.svg)
![Go Version](https://img.shields.io/github/go-mod/go-version/2222-42/bibliography_log)

This directory contains a log of bibliographic references used in the development of Knowledge, especially Zettelkasten.

The system is built using Go and follows Domain-Driven Design (DDD) principles.

## Prerequisites

- Go 1.25.3 or later

## Setup

1. Navigate to the project directory:
   ```bash
   cd path/to/bibliography_log
   ```
2. Download dependencies:
   ```bash
   go mod tidy
   ```

## Usage

The system provides a CLI for managing bibliographies and classifications.

### 1. Add a Classification

Before adding a bibliography, you must ensure the classification exists.

**Command:**
```bash
go run cmd/biblog/*.go add-class -code <code_num> -name "<name>"
```

**Example:**
```bash
go run cmd/biblog/*.go add-class -code 56 -name "Technology"
```
**Output:**
```
Classification added: &{5f450cd1-83e5-49d5-9e67-4becc6ca7efd 56 Technology}
```

### 2. Add a Bibliography

Add a new bibliography entry. The `BibIndex` will be automatically generated based on the input.

**Command:**
```bash
go run cmd/biblog/*.go add-bib -title "<title>" -author "<author>" -type "<type>" -class <class_code> -year <year> -isbn "<isbn>" -desc "<description>"
```

**Example:**
```bash
go run cmd/biblog/*.go add-bib -title "Domain Driven Design" -author "Eric Evans" -type "Book" -class 56 -year 2003 -isbn "978-0321125217" -desc "Tackling Complexity in the Heart of Software"
```
**Output:**
```
Bibliography added: &{f792718c-c789-48b8-8d89-0d8650d4fe35 B56EE03DDD B56 Book Domain Driven Design Eric Evans 978-0321125217 Tackling Complexity in the Heart of Software 2003-01-01 00:00:00 +0000 UTC}
```

### 3. List Bibliographies

List all registered bibliographies.

**Command:**
```bash
go run cmd/biblog/*.go list
```

**Example Output:**
```
Bibliographies:
[Book] Domain Driven Design by Eric Evans (BibIndex: B56EE03DDD)
```

### 4. Add Bibliography with Japanese Text

When adding bibliographies with Japanese titles or authors, you must provide English translations for BibIndex generation.

**Command:**
```bash
go run cmd/biblog/*.go add-bib \
  -title "<japanese_title>" \
  -title-en "<english_title>" \
  -author "<japanese_author>" \
  -author-en "<english_author>" \
  -type "<type>" \
  -class <class_code> \
  -year <year>
```

**Example:**
```bash
go run cmd/biblog/*.go add-bib \
  -title "マネジメント神話" \
  -title-en "The Management Myth" \
  -author "マシュー・スチュワート" \
  -author-en "Matthew Stewart" \
  -type "Book" \
  -class 16 \
  -year 2024 \
  -isbn "978-4750356884"
```

**Output:**
```
Bibliography added: &{<uuid> B16MS24TMM B16 Book マネジメント神話 マシュー・スチュワート 978-4750356884  2024-01-01 00:00:00 +0000 UTC}
```

> **Note:** The system automatically detects Japanese characters (Hiragana, Katakana, Kanji). English translations are only used for generating readable BibIndex codes; the original Japanese text is preserved in the stored data.

### 5. Add Bibliography with Manual BibIndex

You can manually specify the `BibIndex` using the `-bib-index` flag. This overrides the automatic generation logic.

**Command:**
```bash
go run cmd/biblog/*.go add-bib \
  -title "<title>" \
  -author "<author>" \
  -type "<type>" \
  -class <class_code> \
  -year <year> \
  -bib-index "<custom_index>"
```

**Example:**
```bash
go run cmd/biblog/*.go add-bib \
  -title "My Custom Book" \
  -author "John Doe" \
  -type "Book" \
  -class 56 \
  -year 2024 \
  -bib-index "CUSTOM123"
```

**Output:**
```
Bibliography added: &{... BibIndex:CUSTOM123 ...}
```

### 6. Add Review

Add a review for an existing bibliography.

**Command:**
```bash
go run cmd/biblog/*.go add-review \
  -bib-index "<bib_index>" \
  -goals "<goals>" \
  -summary "<summary>"
```

**Example:**
```bash
go run cmd/biblog/*.go add-review \
  -bib-index "B56EE03DDD" \
  -goals "Understand DDD core concepts" \
  -summary "Excellent introduction to the domain layer."
```

**Output:**
```
Review added: &{... Goals:Understand DDD core concepts ...}
```

### 7. Update Review

Update an existing review's goals and/or summary. At least one field must be provided.

**Command:**
```bash
go run cmd/biblog/*.go update-review \
  -review-id "<review-uuid>" \
  -goals "<new-goals>" \
  -summary "<new-summary>"
```

**Example - Update both fields:**
```bash
go run cmd/biblog/*.go update-review \
  -review-id "2d8a26ef-64e4-4718-b913-085fef527d71" \
  -goals "Refined understanding of DDD patterns" \
  -summary "Comprehensive guide covering strategic and tactical design."
```

**Example - Update only summary:**
```bash
go run cmd/biblog/*.go update-review \
  -review-id "2d8a26ef-64e4-4718-b913-085fef527d71" \
  -summary "After reading: Excellent introduction with practical examples."
```

**Output:**
```
Review updated: &{... Goals:Refined understanding of DDD patterns Summary:Comprehensive guide covering strategic and tactical design. UpdatedAt:2025-11-23T04:45:00+09:00}
```

> **Note:** You can find the review UUID from the `data/reviews.csv` file. At least one of `-goals` or `-summary` must be provided. The `UpdatedAt` timestamp is automatically updated.

## Testing

To run the automated tests:

```bash
go test ./internal/...
```

## Data Storage

The data is stored in CSV files in the `data/` directory:
- `data/bibliographies.csv`: Stores bibliography entries.
- `data/classifications.csv`: Stores classification codes.

## Performance Limitations

> **Note:** This system uses CSV files for data storage, which is suitable for small to medium datasets (hundreds to low thousands of entries) but has performance limitations for larger datasets.

**Known Limitations:**

- **Full File Reads:** Methods like `FindByBibIndex()` call `FindAll()`, which reads and parses the entire CSV file on every query. This is inefficient for large datasets.
- **No Indexing:** CSV files don't support indexing, so all searches are O(n) linear scans.
- **Concurrent Access:** The current implementation has potential race conditions when multiple processes access the same CSV file simultaneously (acceptable for single-user CLI usage).

**Recommendations for Production Use:**

- For datasets with **< 1,000 entries**: Current CSV implementation is acceptable
- For datasets with **1,000-10,000 entries**: Consider implementing in-memory caching
- For datasets with **> 10,000 entries**: Migrate to a proper database (SQLite, PostgreSQL, etc.)

The CSV-based approach was chosen for simplicity, portability, and ease of inspection/editing. It's ideal for personal knowledge management and learning DDD principles without database setup overhead.
