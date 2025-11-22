package domain

import (
	"time"

	"github.com/google/uuid"
)

// Review represents a user's evaluation of a bibliography.
type Review struct {
	ID        uuid.UUID
	BookID    uuid.UUID
	Goals     string
	Summary   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

