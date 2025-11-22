package service

import (
	"fmt"
	"strings"
	"time"

	"bibliography_log/internal/domain"

	"github.com/google/uuid"
)

type BibliographyService struct {
	bibRepo   domain.BibliographyRepository
	classRepo domain.BibClassificationRepository
}

func NewBibliographyService(bibRepo domain.BibliographyRepository, classRepo domain.BibClassificationRepository) *BibliographyService {
	return &BibliographyService{
		bibRepo:   bibRepo,
		classRepo: classRepo,
	}
}

func (s *BibliographyService) AddBibliography(title, author, isbn, description, typeStr string, classCodeNum int, publishedDate time.Time) (*domain.Bibliography, error) {
	// 1. Find Classification
	class, err := s.classRepo.FindByCodeNum(classCodeNum)
	if err != nil {
		return nil, fmt.Errorf("failed to find classification: %w", err)
	}
	if class == nil {
		return nil, fmt.Errorf("classification with code %d not found", classCodeNum)
	}

	// 2. Generate BibIndex
	// Format: Code + AuthorInitials + Year + TitleInitials
	// Code: e.g. B56.
	// Wait, the domain model says Code is "B56".
	// But Classification has CodeNum 56 and Name "Technology".
	// We need to construct the Code.
	// The domain model says: Code (String) (e.g., B56("B"(Book)+"56"("Technology)), "E16"("E"(Essay)+"16"("Philosophy")))
	// So we need a way to map Type to a prefix letter.
	// "Book" -> "B", "Essay" -> "E", "Video" -> "V"?
	// I'll assume first letter of Type for now.

	typePrefix := string(typeStr[0])
	code := fmt.Sprintf("%s%d", typePrefix, class.CodeNum)

	authorInitials := generateAuthorInitials(author)
	yearSuffix := publishedDate.Format("06") // Last 2 digits of year
	titleInitials := generateTitleInitials(title)

	bibIndex := fmt.Sprintf("%s%s%s%s", code, authorInitials, yearSuffix, titleInitials)

	// 3. Create Entity
	bib := &domain.Bibliography{
		ID:            uuid.New(),
		BibIndex:      bibIndex,
		Code:          code,
		Type:          typeStr,
		Title:         title,
		Author:        author,
		ISBN:          isbn,
		Description:   description,
		PublishedDate: publishedDate,
	}

	// 4. Save
	if err := s.bibRepo.Save(bib); err != nil {
		return nil, fmt.Errorf("failed to save bibliography: %w", err)
	}

	return bib, nil
}

func (s *BibliographyService) ListBibliographies() ([]*domain.Bibliography, error) {
	return s.bibRepo.FindAll()
}

func (s *BibliographyService) AddClassification(codeNum int, name string) (*domain.BibClassification, error) {
	// Check if classification already exists
	existing, err := s.classRepo.FindByCodeNum(codeNum)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing classification: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("classification with code %d already exists", codeNum)
	}

	class := &domain.BibClassification{
		ID:      uuid.New(),
		CodeNum: codeNum,
		Name:    name,
	}
	if err := s.classRepo.Save(class); err != nil {
		return nil, err
	}
	return class, nil
}

func generateAuthorInitials(author string) string {
	parts := strings.Fields(author)
	if len(parts) == 0 {
		return "XX"
	}
	if len(parts) == 1 {
		if len(parts[0]) >= 2 {
			return strings.ToUpper(parts[0][:2])
		}
		return strings.ToUpper(parts[0]) + "X"
	}
	// First and Last name initials
	first := parts[0]
	last := parts[len(parts)-1]
	return strings.ToUpper(string(first[0]) + string(last[0]))
}

func generateTitleInitials(title string) string {
	// "Book initials up to three letters"
	// Strategy: Take first letter of first 3 words.
	// If fewer than 3 words, take more letters from available words?
	// Let's stick to first letter of words, max 3.
	parts := strings.Fields(title)
	var initials string
	count := 0
	for _, part := range parts {
		// Skip small words? The domain model didn't specify, but usually "The", "A" are skipped.
		// For simplicity, I'll include everything for now or maybe skip common stop words if I want to be fancy.
		// I'll just take all words for now.
		initials += string(part[0])
		count++
		if count >= 3 {
			break
		}
	}
	return strings.ToUpper(initials)
}
