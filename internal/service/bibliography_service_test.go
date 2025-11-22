package service

import (
	"testing"
	"time"

	"bibliography_log/internal/domain"

	"github.com/google/uuid"
)

// MockBibliographyRepository is a mock implementation of domain.BibliographyRepository
type MockBibliographyRepository struct {
	SavedBibliography *domain.Bibliography
	Bibliographies    map[uuid.UUID]*domain.Bibliography
}

func (m *MockBibliographyRepository) Save(b *domain.Bibliography) error {
	m.SavedBibliography = b
	if m.Bibliographies == nil {
		m.Bibliographies = make(map[uuid.UUID]*domain.Bibliography)
	}
	m.Bibliographies[b.ID] = b
	return nil
}

func (m *MockBibliographyRepository) FindAll() ([]*domain.Bibliography, error) {
	var bibs []*domain.Bibliography
	for _, b := range m.Bibliographies {
		bibs = append(bibs, b)
	}
	return bibs, nil
}

func (m *MockBibliographyRepository) FindByID(id uuid.UUID) (*domain.Bibliography, error) {
	if m.Bibliographies == nil {
		return nil, nil
	}
	return m.Bibliographies[id], nil
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

	bib, err := svc.AddBibliography(title, author, isbn, desc, typeStr, classCode, pubDate, "", "", "")

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
	_, err := svc.AddBibliography("", "Author", "ISBN", "Desc", "Book", 56, time.Now(), "", "", "")

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
	_, err := svc.AddBibliography("Title", "", "ISBN", "Desc", "Book", 56, time.Now(), "", "", "")

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
	_, err := svc.AddBibliography("Title", "Author", "ISBN", "Desc", "", 56, time.Now(), "", "", "")

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

func TestContainsJapanese(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"English only", "Hello World", false},
		{"Hiragana", "こんにちは", true},
		{"Katakana", "カタカナ", true},
		{"Kanji", "漢字", true},
		{"Mixed Japanese", "マネジメント神話", true},
		{"Mixed with English", "マシュー・スチュワート(稲岡大志訳)", true},
		{"Numbers and English", "Book 2024", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsJapanese(tt.input)
			if result != tt.expected {
				t.Errorf("containsJapanese(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAddBibliography_JapaneseWithEnglish(t *testing.T) {
	// Setup
	bibRepo := &MockBibliographyRepository{}
	classRepo := &MockBibClassificationRepository{
		Classifications: map[int]*domain.BibClassification{
			16: {CodeNum: 16, Name: "Philosophy"},
		},
	}
	svc := NewBibliographyService(bibRepo, classRepo)

	// Test Case with Japanese title and author, with English translations
	bib, err := svc.AddBibliography(
		"マネジメント神話",
		"マシュー・スチュワート",
		"978-4750356884",
		"",
		"Book",
		16,
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		"The Management Myth",
		"Matthew Stewart",
		"",
	)

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if bib == nil {
		t.Fatal("Expected bibliography to be returned")
	}

	// Verify BibIndex uses English translations
	// Code: B16 (Book + 16)
	// Author: MS (Matthew Stewart)
	// Year: 24 (2024)
	// Title: TMM (The Management Myth)
	expectedBibIndex := "B16MS24TMM"
	if bib.BibIndex != expectedBibIndex {
		t.Errorf("Expected BibIndex %s, got %s", expectedBibIndex, bib.BibIndex)
	}

	// Verify original Japanese is stored
	if bib.Title != "マネジメント神話" {
		t.Errorf("Expected original Japanese title to be stored, got %s", bib.Title)
	}
	if bib.Author != "マシュー・スチュワート" {
		t.Errorf("Expected original Japanese author to be stored, got %s", bib.Author)
	}
}

func TestAddBibliography_JapaneseWithoutEnglish_TitleError(t *testing.T) {
	// Setup
	bibRepo := &MockBibliographyRepository{}
	classRepo := &MockBibClassificationRepository{
		Classifications: map[int]*domain.BibClassification{
			16: {CodeNum: 16, Name: "Philosophy"},
		},
	}
	svc := NewBibliographyService(bibRepo, classRepo)

	// Test Case with Japanese title but no English translation
	_, err := svc.AddBibliography(
		"マネジメント神話",
		"Matthew Stewart",
		"",
		"",
		"Book",
		16,
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		"", // No English title
		"",
		"",
	)

	// Assertions
	if err == nil {
		t.Fatal("Expected error for Japanese title without English translation, got nil")
	}
	if err.Error() != "title contains Japanese characters; please provide English translation via -title-en flag" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestAddBibliography_JapaneseWithoutEnglish_AuthorError(t *testing.T) {
	// Setup
	bibRepo := &MockBibliographyRepository{}
	classRepo := &MockBibClassificationRepository{
		Classifications: map[int]*domain.BibClassification{
			16: {CodeNum: 16, Name: "Philosophy"},
		},
	}
	svc := NewBibliographyService(bibRepo, classRepo)

	// Test Case with Japanese author but no English translation
	_, err := svc.AddBibliography(
		"The Management Myth",
		"マシュー・スチュワート",
		"",
		"",
		"Book",
		16,
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		"",
		"", // No English author
		"",
	)

	// Assertions
	if err == nil {
		t.Fatal("Expected error for Japanese author without English translation, got nil")
	}
	if err.Error() != "author contains Japanese characters; please provide English translation via -author-en flag" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestAddBibliography_ManualBibIndex(t *testing.T) {
	// Setup
	bibRepo := &MockBibliographyRepository{}
	classRepo := &MockBibClassificationRepository{
		Classifications: map[int]*domain.BibClassification{
			56: {CodeNum: 56, Name: "Technology"},
		},
	}
	svc := NewBibliographyService(bibRepo, classRepo)

	// Test Case with manual BibIndex
	manualIndex := "CUSTOM123"
	bib, err := svc.AddBibliography(
		"Domain Driven Design",
		"Eric Evans",
		"978-0321125217",
		"Desc",
		"Book",
		56,
		time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
		"", "", manualIndex,
	)

	// Assertions
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if bib.BibIndex != manualIndex {
		t.Errorf("Expected BibIndex %s, got %s", manualIndex, bib.BibIndex)
	}
}
