package domain

// BibliographyRepository defines the interface for persistence.
type BibliographyRepository interface {
	Save(bibliography *Bibliography) error
	FindAll() ([]*Bibliography, error)
	FindByID(id BibliographyID) (*Bibliography, error)
	FindByBibIndex(bibIndex string) (*Bibliography, error)
}

// ClassificationRepository defines the interface for persistence.
type ClassificationRepository interface {
	Save(classification *Classification) error
	FindAll() ([]*Classification, error)
	FindByCodeNum(codeNum int) (*Classification, error)
}

// ReviewRepository defines the interface for persistence.
type ReviewRepository interface {
	Save(review *Review) error
	FindAll() ([]*Review, error)
	FindByID(id ReviewID) (*Review, error)
	FindByBookID(bookID BibliographyID) ([]*Review, error)
}
