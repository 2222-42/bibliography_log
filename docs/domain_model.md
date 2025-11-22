# Bibliography Log Domain Model

## Ubiquitous Language

- **Bibliography**: A published work of literature, papers or podcasts that can be read and reviewed.
- **Classification**: A classification of a bibliography.
- **Review**: A user's evaluation of a book, consisting of goals and a summary.

## Domain-Specific Types

To improve type safety and prevent accidental misuse of IDs across different entities, the system uses domain-specific ID types instead of raw UUIDs:

- **BibliographyID**: Type-safe wrapper for Bibliography entity IDs
- **ReviewID**: Type-safe wrapper for Review entity IDs  
- **ClassificationID**: Type-safe wrapper for Classification entity IDs

These types provide compile-time safety, preventing errors like passing a ReviewID where a BibliographyID is expected. Each type includes helper methods:
- `String()` - returns UUID string representation
- `UUID()` - returns underlying uuid.UUID
- `NewXXXID()` - generates new random ID
- `ParseXXXID(string)` - parses from string

## Entities

### Bibliography
- **Identity**: `BibliographyID` (UUID)
- **Attributes**:
  - `BibIndex` (String) (e.g., "B56SK24DMD(i.e. B56(Code)+SK(Author's Initials)+24(Bottom two digits of Published Year)+DMD(Book initials up to three letters))")
  - `Code` (String) (e.g., B56("B"(Book)+"56"("Technology)), "E16"("E"(Essay)+"16"("Philosophy")))
  - `Type` (String) (e.g., "Book", "Essay", "Video")
  - `Title` (String)
  - `Author` (String)
  - `ISBN` (String, Value Object)
  - `Description` (String)
  - `PublishedDate` (Date)

> **Note:** `AuthorEn` and `TitleEn` are not attributes of the persisted `Bibliography` entity. They are input parameters used temporarily during BibIndex generation in the service layer and are not stored.

### Classification
- **Identity**: `ClassificationID` (domain-specific type wrapping UUID)
- **Attributes**:
  - `CodeNum` (Integer) (e.g., 56)
  - `Name` (String) (e.g., "Technology")

### Review
- **Identity**: `ReviewID` (domain-specific type wrapping UUID)
- **Attributes**:
  - `BookID` (BibliographyID, Foreign Key)
  - `Goals` (String) - Text field that preserves whitespace and line breaks
  - `Summary` (String) - Text field that preserves whitespace and line breaks
  - `CreatedAt` (DateTime)
  - `UpdatedAt` (DateTime)

> **Note:** Unlike short identifier fields (e.g., `Title`, `Author` in Bibliography which are trimmed), `Goals` and `Summary` are text fields that may contain meaningful whitespace and line breaks. While `TrimSpace()` is used during validation to check for empty content, the actual values are intentionally NOT trimmed during storage to preserve user formatting.

## Aggregates

- **Bibliography Aggregate**: Root is `Bibliography`. Reviews might be considered part of the Book aggregate in some contexts, or separate. For this system, `Review` will be its own aggregate root to allow for independent lifecycle (e.g., a user updating their review without locking the book).

## Services

- **BibliographyService**: Handles book registration and retrieval.
- **BibClassificationService**: Handles classification registration and retrieval.

## Infrastructure

Basically, using CSV files for persistence. The reason of this is that it is flyweight, simple, and easy to change, and I am not decided to use a managed database.

- **BibliographyRepository**: Handles the persistence of `Bibliography` entities.
- **BibClassificationRepository**: Handles the persistence of `BibClassification` entities.
