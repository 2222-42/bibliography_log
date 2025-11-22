# Bibliography Log

This directory contains a log of bibliographic references used in the development of Knowledge, especially Zettelkasten.

The system is built using Go and follows Domain-Driven Design (DDD) principles.

## Prerequisites

- Go 1.23 or later

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

## Testing

To run the automated tests:

```bash
go test ./internal/...
```

## Data Storage

The data is stored in CSV files in the `data/` directory:
- `data/bibliographies.csv`: Stores bibliography entries.
- `data/classifications.csv`: Stores classification codes.
