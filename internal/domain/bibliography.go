package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BibliographyID is a domain-specific type for Bibliography entity IDs.
type BibliographyID uuid.UUID

// String returns the string representation of the BibliographyID.
func (id BibliographyID) String() string {
	return uuid.UUID(id).String()
}

// UUID returns the underlying uuid.UUID value.
func (id BibliographyID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

// NewBibliographyID generates a new random BibliographyID.
func NewBibliographyID() BibliographyID {
	return BibliographyID(uuid.New())
}

// ParseBibliographyID parses a string into a BibliographyID.
func ParseBibliographyID(s string) (BibliographyID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return BibliographyID{}, fmt.Errorf("invalid bibliography ID: %w", err)
	}
	return BibliographyID(id), nil
}

// Bibliography represents a published work.
type Bibliography struct {
	ID            BibliographyID
	BibIndex      string // e.g., "B56SK24DMD"
	Code          string // e.g., "B56"
	Type          string // e.g., "Book", "Essay"
	Title         string
	Author        string
	ISBN          string
	Description   string
	PublishedDate time.Time
}
