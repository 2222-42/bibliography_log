package domain

import "github.com/google/uuid"

// BibClassification represents a classification category.
type BibClassification struct {
	ID      uuid.UUID
	CodeNum int    // e.g., 56
	Name    string // e.g., "Technology"
}

