package infrastructure

import (
	"bibliography_log/internal/domain"
	"os"
	"testing"
	"time"
)

func TestCSVBibliographyRepository_SaveAndFind(t *testing.T) {
	// Setup temporary file
	tmpFile, err := os.CreateTemp("", "bib_test_*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := tmpFile.Close(); err != nil {
			t.Error(err)
		}
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Error(err)
		}
	}()

	repo := NewCSVBibliographyRepository(tmpFile.Name())

	// Test Save
	bib := &domain.Bibliography{
		ID:            domain.NewBibliographyID(),
		BibIndex:      "B56TEST",
		Code:          "B56",
		Type:          "Book",
		Title:         "Test Book",
		Author:        "Test Author",
		Publisher:     "Test Publisher",
		ISBN:          "1234567890",
		PublishedDate: time.Now().Truncate(time.Second), // Truncate to match CSV precision if needed, though RFC3339 handles it well.
	}

	err = repo.Save(bib)
	if err != nil {
		t.Fatalf("Failed to save bibliography: %v", err)
	}

	// Test FindAll (get all records with limit=0, offset=0)
	all, err := repo.FindAll(0, 0)
	if err != nil {
		t.Fatalf("Failed to find all: %v", err)
	}

	if len(all) != 1 {
		t.Errorf("Expected 1 bibliography, got %d", len(all))
	}

	savedBib := all[0]
	if savedBib.ID != bib.ID {
		t.Errorf("Expected ID %v, got %v", bib.ID, savedBib.ID)
	}
	if savedBib.Title != bib.Title {
		t.Errorf("Expected Title %s, got %s", bib.Title, savedBib.Title)
	}

	// Test FindByBibIndex
	found, err := repo.FindByBibIndex("B56TEST")
	if err != nil {
		t.Fatalf("Failed to find by index: %v", err)
	}
	if found == nil {
		t.Fatal("Expected to find bibliography by index")
	}
	if found.ID != bib.ID {
		t.Errorf("Expected found ID %v, got %v", bib.ID, found.ID)
	}
}
