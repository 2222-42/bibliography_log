package domain

import "github.com/google/uuid"

// BibliographyRepository defines the interface for persistence.
type BibliographyRepository interface {
	Save(bibliography *Bibliography) error
	FindAll() ([]*Bibliography, error)
	FindByBibIndex(bibIndex string) (*Bibliography, error)
}

// BibClassificationRepository defines the interface for persistence.
type BibClassificationRepository interface {
	Save(classification *BibClassification) error
	FindAll() ([]*BibClassification, error)
	FindByCodeNum(codeNum int) (*BibClassification, error)
}

// ReviewRepository defines the interface for persistence.
type ReviewRepository interface {
	Save(review *Review) error
	FindAll() ([]*Review, error)
	FindByBookID(bookID uuid.UUID) ([]*Review, error)
}
