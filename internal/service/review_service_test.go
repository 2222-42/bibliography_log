package service

import (
	"bibliography_log/internal/domain"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

// MockReviewRepository for testing
type MockReviewRepository struct {
	Reviews map[uuid.UUID]*domain.Review
}

func (m *MockReviewRepository) Save(review *domain.Review) error {
	if m.Reviews == nil {
		m.Reviews = make(map[uuid.UUID]*domain.Review)
	}
	m.Reviews[review.ID] = review
	return nil
}

func (m *MockReviewRepository) FindAll() ([]*domain.Review, error) {
	var reviews []*domain.Review
	for _, r := range m.Reviews {
		reviews = append(reviews, r)
	}
	return reviews, nil
}

func (m *MockReviewRepository) FindByBookID(bookID uuid.UUID) ([]*domain.Review, error) {
	var reviews []*domain.Review
	for _, r := range m.Reviews {
		if r.BookID == bookID {
			reviews = append(reviews, r)
		}
	}
	return reviews, nil
}

func TestAddReview_Success(t *testing.T) {
	// Setup
	reviewRepo := &MockReviewRepository{}
	bibRepo := &MockBibliographyRepository{
		Bibliographies: map[uuid.UUID]*domain.Bibliography{},
	}

	// Add a dummy bibliography
	bookID := uuid.New()
	bibRepo.Bibliographies[bookID] = &domain.Bibliography{
		ID:    bookID,
		Title: "Test Book",
	}

	svc := NewReviewService(reviewRepo, bibRepo)

	// Test
	goals := "Learn Go"
	summary := "Great book"
	review, err := svc.AddReview(bookID, goals, summary)
	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if review.Goals != goals {
		t.Errorf("Expected goals %s, got %s", goals, review.Goals)
	}
	if review.Summary != summary {
		t.Errorf("Expected summary %s, got %s", summary, review.Summary)
	}
	if review.BookID != bookID {
		t.Errorf("Expected BookID %s, got %s", bookID, review.BookID)
	}
}

func TestAddReview_EmptyGoals(t *testing.T) {
	// Setup
	reviewRepo := &MockReviewRepository{}
	bibRepo := &MockBibliographyRepository{}
	svc := NewReviewService(reviewRepo, bibRepo)

	// Test
	_, err := svc.AddReview(uuid.New(), "", "Summary")

	// Assertions
	if err == nil {
		t.Fatal("Expected error for empty goals, got nil")
	}
	if err.Error() != "goals are required and cannot be empty" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestAddReview_WhitespaceGoals(t *testing.T) {
	// Setup
	reviewRepo := &MockReviewRepository{}
	bibRepo := &MockBibliographyRepository{}
	svc := NewReviewService(reviewRepo, bibRepo)

	// Test Case with whitespace-only goals
	_, err := svc.AddReview(uuid.New(), "   ", "Summary")

	// Assertions
	if err == nil {
		t.Fatal("Expected error for whitespace-only goals, got nil")
	}
	if err.Error() != "goals are required and cannot be empty" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestAddReview_BookNotFound(t *testing.T) {
	// Setup
	reviewRepo := &MockReviewRepository{}
	bibRepo := &MockBibliographyRepository{
		Bibliographies: map[uuid.UUID]*domain.Bibliography{},
	}
	svc := NewReviewService(reviewRepo, bibRepo)

	// Test
	nonExistentID := uuid.New()
	_, err := svc.AddReview(nonExistentID, "Goals", "Summary")

	// Assertions
	if err == nil {
		t.Fatal("Expected error for non-existent book, got nil")
	}
	expectedErr := fmt.Sprintf("bibliography with ID %s not found", nonExistentID)
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}
