# Bibliography Log Domain Model

## Ubiquitous Language

- **Bibliography**: A published work of literature, papers or podcasts that can be read and reviewed.
- **BibClassification**: A classification of a bibliography.
- **Review**: A user's evaluation of a book, consisting of a rating and an optional comment.

## Entities

### Bibliography
- **Identity**: `BookID` (UUID)
- **Attributes**:
  - `BibIndex` (String) (e.g., "B56SK24DMD(i.e. B56(Code)+SK(Author's Initials)+24(Bottom two digits of Published Year)+DMD(Book initials up to three letters))")
  - `Code` (String) (e.g., B56("B"(Book)+"56"("Technology)), "E16"("E"(Essay)+"16"("Philosophy")))
  - `Type` (String) (e.g., "Book", "Essay", "Video")
  - `Title` (String)
  - `Author` (String)
  - `ISBN` (String, Value Object)
  - `Description` (String)
  - `PublishedDate` (Date)

### BibClassification
- **Identity**: `BibClassificationID` (UUID)
- **Attributes**:
  - `CodeNum` (Integer) (e.g., 56)
  - `Name` (String) (e.g., "Technology")

### Review
- **Identity**: `ReviewID` (UUID)
- **Attributes**:
  - `BookID` (UUID, Foreign Key)
  - `Goals` (String)
  - `Summary` (String)
  - `CreatedAt` (DateTime)
  - `UpdatedAt` (DateTime)

## Aggregates

- **Bibliography Aggregate**: Root is `Bibliography`. Reviews might be considered part of the Book aggregate in some contexts, or separate. For this system, `Review` will be its own aggregate root to allow for independent lifecycle (e.g., a user updating their review without locking the book).

## Services

- **BibliographyService**: Handles book registration and retrieval.
- **BibClassificationService**: Handles classification registration and retrieval.

## Infrastructure

Basically, using CSV files for persistence. The reason of this is that it is flyweight, simple, and easy to change, and I am not decided to use a managed database.

- **BibliographyRepository**: Handles the persistence of `Bibliography` entities.
- **BibClassificationRepository**: Handles the persistence of `BibClassification` entities.
