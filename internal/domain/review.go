package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ReviewID is a domain-specific type for Review entity IDs.
type ReviewID uuid.UUID

// String returns the string representation of the ReviewID.
func (id ReviewID) String() string {
	return uuid.UUID(id).String()
}

// UUID returns the underlying uuid.UUID value.
func (id ReviewID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

// NewReviewID generates a new random ReviewID.
func NewReviewID() ReviewID {
	return ReviewID(uuid.New())
}

// ParseReviewID parses a string into a ReviewID.
func ParseReviewID(s string) (ReviewID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ReviewID{}, fmt.Errorf("invalid review ID: %w", err)
	}
	return ReviewID(id), nil
}

// Review represents a user's evaluation of a bibliography.
type Review struct {
	ID        ReviewID
	BookID    BibliographyID
	Goals     string
	Summary   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
