package service

import (
	"testing"
	"time"

	"bibliography_log/internal/domain"
)

// MockBibliographyRepository is a mock implementation of domain.BibliographyRepository
type MockBibliographyRepository struct {
	SavedBibliography *domain.Bibliography
}

func (m *MockBibliographyRepository) Save(b *domain.Bibliography) error {
	m.SavedBibliography = b
	return nil
}

func (m *MockBibliographyRepository) FindAll() ([]*domain.Bibliography, error) {
	return nil, nil
}

func (m *MockBibliographyRepository) FindByBibIndex(bibIndex string) (*domain.Bibliography, error) {
	return nil, nil
}

// MockBibClassificationRepository is a mock implementation of domain.BibClassificationRepository
type MockBibClassificationRepository struct {
	Classifications map[int]*domain.BibClassification
}

func (m *MockBibClassificationRepository) Save(c *domain.BibClassification) error {
	if m.Classifications == nil {
		m.Classifications = make(map[int]*domain.BibClassification)
	}
	m.Classifications[c.CodeNum] = c
	return nil
}

func (m *MockBibClassificationRepository) FindAll() ([]*domain.BibClassification, error) {
	return nil, nil
}

func (m *MockBibClassificationRepository) FindByCodeNum(codeNum int) (*domain.BibClassification, error) {
	if c, ok := m.Classifications[codeNum]; ok {
		return c, nil
	}
	return nil, nil
}

func TestAddBibliography(t *testing.T) {
	// Setup
	bibRepo := &MockBibliographyRepository{}
	classRepo := &MockBibClassificationRepository{
		Classifications: map[int]*domain.BibClassification{
			56: {CodeNum: 56, Name: "Technology"},
		},
	}
	svc := NewBibliographyService(bibRepo, classRepo)

	// Test Case
	title := "Domain Driven Design"
	author := "Eric Evans"
	isbn := "978-0321125217"
	desc := "Tackling Complexity"
	typeStr := "Book"
	classCode := 56
	pubDate := time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC)

	bib, err := svc.AddBibliography(title, author, isbn, desc, typeStr, classCode, pubDate)

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if bib == nil {
		t.Fatal("Expected bibliography to be returned")
	}

	// Verify BibIndex generation
	// Code: B56 (First letter of Book + 56)
	// Author: EE (Eric Evans)
	// Year: 03 (2003)
	// Title: DDD (Domain Driven Design)
	expectedBibIndex := "B56EE03DDD"
	if bib.BibIndex != expectedBibIndex {
		t.Errorf("Expected BibIndex %s, got %s", expectedBibIndex, bib.BibIndex)
	}

	if bibRepo.SavedBibliography != bib {
		t.Error("Expected bibliography to be saved to repository")
	}
}

func TestAddClassification(t *testing.T) {
	// Setup
	bibRepo := &MockBibliographyRepository{}
	classRepo := &MockBibClassificationRepository{}
	svc := NewBibliographyService(bibRepo, classRepo)

	// Test Case
	codeNum := 99
	name := "Test Class"

	class, err := svc.AddClassification(codeNum, name)

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if class == nil {
		t.Fatal("Expected classification to be returned")
	}

	if class.CodeNum != codeNum {
		t.Errorf("Expected CodeNum %d, got %d", codeNum, class.CodeNum)
	}

	if class.Name != name {
		t.Errorf("Expected Name %s, got %s", name, class.Name)
	}

	// Verify it was saved
	saved, _ := classRepo.FindByCodeNum(codeNum)
	if saved == nil {
		t.Error("Expected classification to be saved to repository")
	}
}

func TestAddClassification_Duplicate(t *testing.T) {
	// Setup
	bibRepo := &MockBibliographyRepository{}
	classRepo := &MockBibClassificationRepository{
		Classifications: map[int]*domain.BibClassification{
			99: {CodeNum: 99, Name: "Existing Class"},
		},
	}
	svc := NewBibliographyService(bibRepo, classRepo)

	// Test Case
	codeNum := 99
	name := "New Class"

	_, err := svc.AddClassification(codeNum, name)

	// Assertions
	if err == nil {
		t.Fatal("Expected error for duplicate classification, got nil")
	}
}

func TestAddBibliography_EmptyTitle(t *testing.T) {
	// Setup
	bibRepo := &MockBibliographyRepository{}
	classRepo := &MockBibClassificationRepository{
		Classifications: map[int]*domain.BibClassification{
			56: {CodeNum: 56, Name: "Technology"},
		},
	}
	svc := NewBibliographyService(bibRepo, classRepo)

	// Test Case with empty title
	_, err := svc.AddBibliography("", "Author", "ISBN", "Desc", "Book", 56, time.Now())

	// Assertions
	if err == nil {
		t.Fatal("Expected error for empty title, got nil")
	}
	if err.Error() != "title is required and cannot be empty" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestAddBibliography_EmptyAuthor(t *testing.T) {
	// Setup
	bibRepo := &MockBibliographyRepository{}
	classRepo := &MockBibClassificationRepository{
		Classifications: map[int]*domain.BibClassification{
			56: {CodeNum: 56, Name: "Technology"},
		},
	}
	svc := NewBibliographyService(bibRepo, classRepo)

	// Test Case with empty author
	_, err := svc.AddBibliography("Title", "", "ISBN", "Desc", "Book", 56, time.Now())

	// Assertions
	if err == nil {
		t.Fatal("Expected error for empty author, got nil")
	}
	if err.Error() != "author is required and cannot be empty" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestAddBibliography_EmptyType(t *testing.T) {
	// Setup
	bibRepo := &MockBibliographyRepository{}
	classRepo := &MockBibClassificationRepository{
		Classifications: map[int]*domain.BibClassification{
			56: {CodeNum: 56, Name: "Technology"},
		},
	}
	svc := NewBibliographyService(bibRepo, classRepo)

	// Test Case with empty type
	_, err := svc.AddBibliography("Title", "Author", "ISBN", "Desc", "", 56, time.Now())

	// Assertions
	if err == nil {
		t.Fatal("Expected error for empty type, got nil")
	}
	if err.Error() != "type is required and cannot be empty" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestAddClassification_EmptyName(t *testing.T) {
	// Setup
	bibRepo := &MockBibliographyRepository{}
	classRepo := &MockBibClassificationRepository{}
	svc := NewBibliographyService(bibRepo, classRepo)

	// Test Case with empty name
	_, err := svc.AddClassification(99, "")

	// Assertions
	if err == nil {
		t.Fatal("Expected error for empty name, got nil")
	}
	if err.Error() != "classification name must not be empty" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestAddClassification_WhitespaceName(t *testing.T) {
	// Setup
	bibRepo := &MockBibliographyRepository{}
	classRepo := &MockBibClassificationRepository{}
	svc := NewBibliographyService(bibRepo, classRepo)

	// Test Case with whitespace-only name
	_, err := svc.AddClassification(99, "   ")

	// Assertions
	if err == nil {
		t.Fatal("Expected error for whitespace-only name, got nil")
	}
	if err.Error() != "classification name must not be empty" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}
