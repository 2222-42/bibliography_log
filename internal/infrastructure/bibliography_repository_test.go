package infrastructure

import (
	"bibliography_log/internal/domain"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCSVBibliographyRepository_SaveAndFind(t *testing.T) {
	// Setup temporary file
	tmpFile, err := os.CreateTemp("", "bib_test_*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile.Close()           // ensure file handle is closed
	defer os.Remove(tmpFile.Name()) // clean up

	repo := NewCSVBibliographyRepository(tmpFile.Name())

	// Test Save
	bib := &domain.Bibliography{
		ID:            uuid.New(),
		BibIndex:      "B56TEST",
		Code:          "B56",
		Type:          "Book",
		Title:         "Test Book",
		Author:        "Test Author",
		ISBN:          "1234567890",
		Description:   "Test Description",
		PublishedDate: time.Now().Truncate(time.Second), // Truncate to match CSV precision if needed, though RFC3339 handles it well.
	}

	err = repo.Save(bib)
	if err != nil {
		t.Fatalf("Failed to save bibliography: %v", err)
	}

	// Test FindAll
	all, err := repo.FindAll()
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
