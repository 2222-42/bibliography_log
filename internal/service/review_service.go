package service

import (
	"bibliography_log/internal/domain"
	"fmt"
	"time"

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
	if goals == "" {
		return nil, fmt.Errorf("goals are required and cannot be empty")
	}

	// Verify book exists
	// Note: Ideally we should have FindByID, but for now we rely on the caller to provide a valid ID
	// or we can implement FindByID in BibliographyRepository if needed.
	// Given the current repository interface, we might need to rely on the caller or add FindByID.
	// Let's assume for now that the caller (CLI) resolves the ID correctly.
	// However, to be safe, we should verify existence if possible.
	// Since FindAll is available, we can use it to verify.
	allBibs, err := s.bibRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to verify book existence: %w", err)
	}
	found := false
	for _, b := range allBibs {
		if b.ID == bookID {
			found = true
			break
		}
	}
	if !found {
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
