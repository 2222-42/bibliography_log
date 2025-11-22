package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// ClassificationID is a domain-specific type for Classification entity IDs.
type ClassificationID uuid.UUID

// String returns the string representation of the ClassificationID.
func (id ClassificationID) String() string {
	return uuid.UUID(id).String()
}

// UUID returns the underlying uuid.UUID value.
func (id ClassificationID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

// NewClassificationID generates a new random ClassificationID.
func NewClassificationID() ClassificationID {
	return ClassificationID(uuid.New())
}

// ParseClassificationID parses a string into a ClassificationID.
func ParseClassificationID(s string) (ClassificationID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ClassificationID{}, fmt.Errorf("invalid classification ID: %w", err)
	}
	return ClassificationID(id), nil
}

// Classification represents a classification category.
type Classification struct {
	ID      ClassificationID
	CodeNum int    // e.g., 56
	Name    string // e.g., "Technology"
}
