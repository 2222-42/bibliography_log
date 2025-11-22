package service

import (
	"fmt"
	"strings"
	"time"

	"bibliography_log/internal/domain"

	"github.com/google/uuid"
)

type ReviewService struct {
	reviewRepo domain.ReviewRepository
	bibRepo    domain.BibliographyRepository
}

func NewReviewService(reviewRepo domain.ReviewRepository, bibRepo domain.BibliographyRepository) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		bibRepo:    bibRepo,
	}
}

func (s *ReviewService) AddReview(bookID uuid.UUID, goals string, summary string) (*domain.Review, error) {
	// Validate inputs
	// Note: 'goals' and 'summary' are text fields that may contain meaningful whitespace
	// and line breaks, so we do NOT trim them before storage (unlike short identifier fields
	// like 'title' or 'author' which are trimmed in BibliographyService.AddBibliography).
	// We only use TrimSpace() for validation to check if the content is non-empty.
	// Note: 'summary' is optional and does not require validation. If this changes, add validation here.
	if strings.TrimSpace(goals) == "" {
		return nil, fmt.Errorf("goals are required and cannot be empty")
	}

	// Verify book exists
	// Use FindByID for efficient existence check.
	bib, err := s.bibRepo.FindByID(bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify book existence: %w", err)
	}
	if bib == nil {
		return nil, fmt.Errorf("bibliography with ID %s not found", bookID)
	}

	review := &domain.Review{
		ID:        uuid.New(),
		BookID:    bookID,
		Goals:     goals,
		Summary:   summary,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.reviewRepo.Save(review); err != nil {
		return nil, fmt.Errorf("failed to save review: %w", err)
	}

	return review, nil
}
