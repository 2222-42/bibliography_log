package domain

import (
	"time"

	"github.com/google/uuid"
)

// Bibliography represents a published work.
type Bibliography struct {
	ID            uuid.UUID
	BibIndex      string // e.g., "B56SK24DMD"
	Code          string // e.g., "B56"
	Type          string // e.g., "Book", "Essay"
	Title         string
	Author        string
	ISBN          string
	Description   string
	PublishedDate time.Time
}
