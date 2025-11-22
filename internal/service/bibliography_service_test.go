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
