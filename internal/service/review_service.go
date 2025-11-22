package service

import (
	"bibliography_log/internal/domain"
	"fmt"
	"strings"
	"time"
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

func (s *ReviewService) AddReview(bookID domain.BibliographyID, goals string, summary string) (*domain.Review, error) {
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
		ID:        domain.NewReviewID(),
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

// UpdateReview updates an existing review's goals and/or summary.
// At least one of goals or summary must be provided (non-nil pointer).
// If a field is nil, it will not be updated (preserves existing value).
// For goals: if provided, must be non-empty/non-whitespace (cannot be set to empty string).
// For summary: if provided, can be set to empty string (no validation).
func (s *ReviewService) UpdateReview(id domain.ReviewID, goals *string, summary *string) (*domain.Review, error) {
	// Validate that at least one field is being updated
	if goals == nil && summary == nil {
		return nil, fmt.Errorf("at least one field (goals or summary) must be provided for update")
	}

	// Retrieve existing review
	review, err := s.reviewRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find review: %w", err)
	}
	if review == nil {
		return nil, fmt.Errorf("review with ID %s not found", id)
	}

	// Update fields if provided
	if goals != nil {
		// Validate goals if being updated (same validation as AddReview)
		if strings.TrimSpace(*goals) == "" {
			return nil, fmt.Errorf("goals cannot be empty or whitespace-only")
		}
		review.Goals = *goals
	}
	if summary != nil {
		// Summary can be empty, so no validation needed
		review.Summary = *summary
	}

	// Update timestamp
	review.UpdatedAt = time.Now()

	// Save the updated review
	if err := s.reviewRepo.Save(review); err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	return review, nil
}
