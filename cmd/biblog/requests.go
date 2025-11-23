package main

import (
	"bibliography_log/internal/domain"
	"fmt"
	"time"
)

// AddClassificationRequest holds arguments for adding a classification.
type AddClassificationRequest struct {
	Code int
	Name string
}

func (r *AddClassificationRequest) PromptMissing() {
	if r.Code == 0 {
		r.Code = promptInt("Classification Code Number", true)
	}
	if r.Name == "" {
		r.Name = promptString("Classification Name", true)
	}
}

func (r *AddClassificationRequest) Validate() error {
	if r.Code == 0 {
		return fmt.Errorf("classification code is required")
	}
	if r.Name == "" {
		return fmt.Errorf("classification name is required")
	}
	return nil
}

// AddBibliographyRequest holds arguments for adding a bibliography.
type AddBibliographyRequest struct {
	Title     string
	Author    string
	Publisher string
	Type      string
	ClassCode int
	Year      int
	ISBN      string
	TitleEn   string
	AuthorEn  string
	BibIndex  string
}

func (r *AddBibliographyRequest) PromptMissing() {
	if r.Title == "" {
		r.Title = promptString("Title", true)
	}
	if r.Author == "" {
		r.Author = promptString("Author", true)
	}
	if r.Publisher == "" {
		r.Publisher = promptString("Publisher", false)
	}
	if r.Type == "" {
		r.Type = promptString("Type", true)
	}
	if r.ClassCode == 0 {
		r.ClassCode = promptInt("Classification Code Number", true)
	}
	if r.Year == 0 {
		r.Year = promptInt("Published Year", true)
	}
	if r.ISBN == "" {
		r.ISBN = promptString("ISBN", false)
	}
	if r.BibIndex == "" {
		r.BibIndex = promptString("BibIndex", false)
	}
}

func (r *AddBibliographyRequest) Validate() error {
	if r.Title == "" {
		return fmt.Errorf("title is required")
	}
	if r.Author == "" {
		return fmt.Errorf("author is required")
	}
	if r.Type == "" {
		return fmt.Errorf("type is required")
	}
	if r.ClassCode == 0 {
		return fmt.Errorf("classification code is required")
	}
	if r.Year == 0 {
		return fmt.Errorf("published year is required")
	}
	return nil
}

func (r *AddBibliographyRequest) ToPublishedDate() time.Time {
	return time.Date(r.Year, 1, 1, 0, 0, 0, 0, time.UTC)
}

// AddReviewRequest holds arguments for adding a review.
type AddReviewRequest struct {
	BibIndex string
	Goals    string
	Summary  string
}

func (r *AddReviewRequest) PromptMissing() {
	if r.BibIndex == "" {
		r.BibIndex = promptString("BibIndex", true)
	}
	if r.Goals == "" {
		r.Goals = promptString("Goals", true)
	}
	if r.Summary == "" {
		r.Summary = promptString("Summary", false)
	}
}

func (r *AddReviewRequest) Validate() error {
	if r.BibIndex == "" {
		return fmt.Errorf("bib-index is required")
	}
	if r.Goals == "" {
		return fmt.Errorf("goals are required")
	}
	return nil
}

// UpdateReviewRequest holds arguments for updating a review.
type UpdateReviewRequest struct {
	ReviewIDStr string
	Goals       string
	Summary     string
}

func (r *UpdateReviewRequest) PromptMissing() {
	if r.ReviewIDStr == "" {
		r.ReviewIDStr = promptString("Review UUID", true)
	}
	// If neither optional field is provided, prompt for them
	if r.Goals == "" && r.Summary == "" {
		r.Goals = promptString("New goals", false)
		r.Summary = promptString("New summary", false)
	}
}

func (r *UpdateReviewRequest) Validate() error {
	if r.ReviewIDStr == "" {
		return fmt.Errorf("review-id is required")
	}
	if r.Goals == "" && r.Summary == "" {
		return fmt.Errorf("at least one field to update (goals or summary) is required")
	}
	return nil
}

func (r *UpdateReviewRequest) ParseID() (domain.ReviewID, error) {
	return domain.ParseReviewID(r.ReviewIDStr)
}

// ListBibliographiesRequest holds arguments for listing bibliographies.
type ListBibliographiesRequest struct {
	Limit  int
	Offset int
}

func (r *ListBibliographiesRequest) Validate() error {
	if r.Limit < 0 {
		return fmt.Errorf("limit must be non-negative")
	}
	if r.Offset < 0 {
		return fmt.Errorf("offset must be non-negative")
	}
	return nil
}
