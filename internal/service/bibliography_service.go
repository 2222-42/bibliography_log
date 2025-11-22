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
	// Input validation for required fields
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("title is required and cannot be empty")
	}
	if strings.TrimSpace(author) == "" {
		return nil, fmt.Errorf("author is required and cannot be empty")
	}
	if strings.TrimSpace(typeStr) == "" {
		return nil, fmt.Errorf("type is required and cannot be empty")
	}

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
	// The Code is constructed by concatenating a type prefix (first letter of the type string, e.g. "B" for "Book")
	// with the classification code number (e.g. 56 for "Technology").
	// Example: "Book" type and classification code 56 yields "B56".
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
	// Validate name is not empty or whitespace
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("classification name must not be empty")
	}

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
		return nil, fmt.Errorf("failed to save classification: %w", err)
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
	var firstInitial, lastInitial string
	if len(first) > 0 {
		firstInitial = string(first[0])
	} else {
		firstInitial = "X"
	}
	if len(last) > 0 {
		lastInitial = string(last[0])
	} else {
		lastInitial = "X"
	}
	return strings.ToUpper(firstInitial + lastInitial)
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
		// Note: All words are currently included when generating initials, including small words like "The" and "A".
		// Skipping common stop words is not implemented, but could be added as a future enhancement.
		// The current implementation takes the first letter of the first three words, regardless of word size.
		if len(part) == 0 {
			continue
		}
		initials += string(part[0])
		count++
		if count >= 3 {
			break
		}
	}
	return strings.ToUpper(initials)
}
