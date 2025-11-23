package service

import (
	"bibliography_log/internal/domain"
	"fmt"
	"testing"
	"time"
)

// MockReviewRepository for testing
type MockReviewRepository struct {
	Reviews map[domain.ReviewID]*domain.Review
}

func (m *MockReviewRepository) Save(review *domain.Review) error {
	if m.Reviews == nil {
		m.Reviews = make(map[domain.ReviewID]*domain.Review)
	}
	m.Reviews[review.ID] = review
	return nil
}

func (m *MockReviewRepository) FindAll(limit, offset int) ([]*domain.Review, error) {
	var reviews []*domain.Review
	for _, r := range m.Reviews {
		reviews = append(reviews, r)
	}
	// Apply offset
	if offset > 0 {
		if offset >= len(reviews) {
			return []*domain.Review{}, nil
		}
		reviews = reviews[offset:]
	}
	// Apply limit
	if limit > 0 && limit < len(reviews) {
		reviews = reviews[:limit]
	}
	return reviews, nil
}

func (m *MockReviewRepository) FindByID(id domain.ReviewID) (*domain.Review, error) {
	if m.Reviews == nil {
		return nil, nil
	}
	return m.Reviews[id], nil
}

func (m *MockReviewRepository) FindByBookID(bookID domain.BibliographyID) ([]*domain.Review, error) {
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
		Bibliographies: map[domain.BibliographyID]*domain.Bibliography{},
	}

	// Add a dummy bibliography
	bookID := domain.NewBibliographyID()
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
	_, err := svc.AddReview(domain.NewBibliographyID(), "", "Summary")

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
	_, err := svc.AddReview(domain.NewBibliographyID(), "   ", "Summary")

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
		Bibliographies: map[domain.BibliographyID]*domain.Bibliography{},
	}
	svc := NewReviewService(reviewRepo, bibRepo)

	// Test
	nonExistentID := domain.NewBibliographyID()
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

func TestUpdateReview_Success(t *testing.T) {
	// Setup
	reviewRepo := &MockReviewRepository{
		Reviews: make(map[domain.ReviewID]*domain.Review),
	}
	bibRepo := &MockBibliographyRepository{}

	// Create an initial review
	reviewID := domain.NewReviewID()
	bookID := domain.NewBibliographyID()
	initialReview := &domain.Review{
		ID:        reviewID,
		BookID:    bookID,
		Goals:     "Initial goals",
		Summary:   "Initial summary",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	}
	reviewRepo.Reviews[reviewID] = initialReview

	svc := NewReviewService(reviewRepo, bibRepo)

	// Test updating both fields
	newGoals := "Updated goals"
	newSummary := "Updated summary"
	updated, err := svc.UpdateReview(reviewID, &newGoals, &newSummary)

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if updated.Goals != newGoals {
		t.Errorf("Expected goals '%s', got '%s'", newGoals, updated.Goals)
	}
	if updated.Summary != newSummary {
		t.Errorf("Expected summary '%s', got '%s'", newSummary, updated.Summary)
	}
	// UpdatedAt should be updated (greater than or equal to initial, since timing can be very fast)
	if updated.UpdatedAt.Before(initialReview.UpdatedAt) {
		t.Errorf("Expected UpdatedAt to be updated or equal, but it's before the initial time")
	}
	if !updated.CreatedAt.Equal(initialReview.CreatedAt) {
		t.Errorf("Expected CreatedAt to remain unchanged")
	}
}

func TestUpdateReview_OnlyGoals(t *testing.T) {
	// Setup
	reviewRepo := &MockReviewRepository{
		Reviews: make(map[domain.ReviewID]*domain.Review),
	}
	bibRepo := &MockBibliographyRepository{}

	reviewID := domain.NewReviewID()
	initialReview := &domain.Review{
		ID:      reviewID,
		BookID:  domain.NewBibliographyID(),
		Goals:   "Initial goals",
		Summary: "Initial summary",
	}
	reviewRepo.Reviews[reviewID] = initialReview

	svc := NewReviewService(reviewRepo, bibRepo)

	// Test updating only goals
	newGoals := "Updated goals only"
	updated, err := svc.UpdateReview(reviewID, &newGoals, nil)

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if updated.Goals != newGoals {
		t.Errorf("Expected goals '%s', got '%s'", newGoals, updated.Goals)
	}
	if updated.Summary != initialReview.Summary {
		t.Errorf("Expected summary to remain unchanged as '%s', got '%s'", initialReview.Summary, updated.Summary)
	}
}

func TestUpdateReview_OnlySummary(t *testing.T) {
	// Setup
	reviewRepo := &MockReviewRepository{
		Reviews: make(map[domain.ReviewID]*domain.Review),
	}
	bibRepo := &MockBibliographyRepository{}

	reviewID := domain.NewReviewID()
	initialReview := &domain.Review{
		ID:      reviewID,
		BookID:  domain.NewBibliographyID(),
		Goals:   "Initial goals",
		Summary: "Initial summary",
	}
	reviewRepo.Reviews[reviewID] = initialReview

	svc := NewReviewService(reviewRepo, bibRepo)

	// Test updating only summary
	newSummary := "Updated summary only"
	updated, err := svc.UpdateReview(reviewID, nil, &newSummary)

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if updated.Goals != initialReview.Goals {
		t.Errorf("Expected goals to remain unchanged as '%s', got '%s'", initialReview.Goals, updated.Goals)
	}
	if updated.Summary != newSummary {
		t.Errorf("Expected summary '%s', got '%s'", newSummary, updated.Summary)
	}
}

func TestUpdateReview_NotFound(t *testing.T) {
	// Setup
	reviewRepo := &MockReviewRepository{
		Reviews: make(map[domain.ReviewID]*domain.Review),
	}
	bibRepo := &MockBibliographyRepository{}
	svc := NewReviewService(reviewRepo, bibRepo)

	// Test updating non-existent review
	nonExistentID := domain.NewReviewID()
	newGoals := "Some goals"
	_, err := svc.UpdateReview(nonExistentID, &newGoals, nil)

	// Assertions
	if err == nil {
		t.Fatal("Expected error for non-existent review, got nil")
	}
	expectedErr := fmt.Sprintf("review with ID %s not found", nonExistentID)
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestUpdateReview_NoFieldsProvided(t *testing.T) {
	// Setup
	reviewRepo := &MockReviewRepository{}
	bibRepo := &MockBibliographyRepository{}
	svc := NewReviewService(reviewRepo, bibRepo)

	// Test updating without providing any fields
	_, err := svc.UpdateReview(domain.NewReviewID(), nil, nil)

	// Assertions
	if err == nil {
		t.Fatal("Expected error when no fields provided, got nil")
	}
	if err.Error() != "at least one field (goals or summary) must be provided for update" {
		t.Errorf("Expected specific error message, got '%s'", err.Error())
	}
}

func TestUpdateReview_EmptyGoals(t *testing.T) {
	// Setup
	reviewRepo := &MockReviewRepository{
		Reviews: make(map[domain.ReviewID]*domain.Review),
	}
	bibRepo := &MockBibliographyRepository{}

	reviewID := domain.NewReviewID()
	initialReview := &domain.Review{
		ID:      reviewID,
		BookID:  domain.NewBibliographyID(),
		Goals:   "Initial goals",
		Summary: "Initial summary",
	}
	reviewRepo.Reviews[reviewID] = initialReview

	svc := NewReviewService(reviewRepo, bibRepo)

	// Test updating with empty goals
	emptyGoals := ""
	_, err := svc.UpdateReview(reviewID, &emptyGoals, nil)

	// Assertions
	if err == nil {
		t.Fatal("Expected error for empty goals, got nil")
	}
	if err.Error() != "goals cannot be empty or whitespace-only" {
		t.Errorf("Expected specific error message, got '%s'", err.Error())
	}
}
