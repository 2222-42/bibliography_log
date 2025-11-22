package infrastructure

import (
	"bibliography_log/internal/domain"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// CSVBibliographyRepository implements domain.BibliographyRepository using a CSV file.
type CSVBibliographyRepository struct {
	FilePath string
}

func NewCSVBibliographyRepository(filePath string) *CSVBibliographyRepository {
	return &CSVBibliographyRepository{FilePath: filePath}
}

func (r *CSVBibliographyRepository) Save(b *domain.Bibliography) error {
	all, err := r.FindAll()
	if err != nil {
		return err
	}

	updated := false
	for i, existing := range all {
		if existing.ID == b.ID {
			all[i] = b
			updated = true
			break
		}
	}
	if !updated {
		all = append(all, b)
	}

	return r.writeAll(all)
}

func (r *CSVBibliographyRepository) FindAll() ([]*domain.Bibliography, error) {
	records, err := ReadCSV(r.FilePath)
	if err != nil {
		return nil, err
	}

	var bibliographies []*domain.Bibliography
	// Skip header
	if len(records) > 0 {
		records = records[1:]
	}

	for _, record := range records {
		if len(record) < 9 {
			continue
		}
		id, err := uuid.Parse(record[0])
		if err != nil {
			slog.Error("Failed to parse published date", "err", err)
			continue
		}
		pubDate, err := time.Parse(time.RFC3339, record[8])
		if err != nil {
			slog.Error("Failed to parse published date", "err", err)
			continue
		}

		bibliographies = append(bibliographies, &domain.Bibliography{
			ID:            id,
			BibIndex:      record[1],
			Code:          record[2],
			Type:          record[3],
			Title:         record[4],
			Author:        record[5],
			ISBN:          record[6],
			Description:   record[7],
			PublishedDate: pubDate,
		})
	}
	return bibliographies, nil
}

func (r *CSVBibliographyRepository) FindByBibIndex(bibIndex string) (*domain.Bibliography, error) {
	all, err := r.FindAll()
	if err != nil {
		return nil, err
	}
	for _, b := range all {
		if b.BibIndex == bibIndex {
			return b, nil
		}
	}
	return nil, nil
}

func (r *CSVBibliographyRepository) writeAll(bibliographies []*domain.Bibliography) error {
	var records [][]string
	// Header
	records = append(records, []string{"ID", "BibIndex", "Code", "Type", "Title", "Author", "ISBN", "Description", "PublishedDate"})

	for _, b := range bibliographies {
		record := []string{
			b.ID.String(),
			b.BibIndex,
			b.Code,
			b.Type,
			b.Title,
			b.Author,
			b.ISBN,
			b.Description,
			b.PublishedDate.Format(time.RFC3339),
		}
		records = append(records, record)
	}
	return WriteCSV(r.FilePath, records)
}
