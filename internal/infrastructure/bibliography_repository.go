package infrastructure

import (
	"bibliography_log/internal/domain"
	"fmt"
	"log/slog"
	"time"
)

// BibliographyRecord represents a bibliography record for CSV persistence.
type BibliographyRecord struct {
	ID            string
	BibIndex      string
	Code          string
	Type          string
	Title         string
	Author        string
	Publisher     string
	ISBN          string
	PublishedDate string
}

// recordToBibliography converts a BibliographyRecord to a domain.Bibliography.
func recordToBibliography(rec *BibliographyRecord) (*domain.Bibliography, error) {
	id, err := domain.ParseBibliographyID(rec.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bibliography ID: %w", err)
	}

	pubDate, err := time.Parse(time.RFC3339, rec.PublishedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse published date: %w", err)
	}

	return &domain.Bibliography{
		ID:            id,
		BibIndex:      rec.BibIndex,
		Code:          rec.Code,
		Type:          rec.Type,
		Title:         rec.Title,
		Author:        rec.Author,
		Publisher:     rec.Publisher,
		ISBN:          rec.ISBN,
		PublishedDate: pubDate,
	}, nil
}

// bibliographyToRecord converts a domain.Bibliography to a BibliographyRecord.
func bibliographyToRecord(bib *domain.Bibliography) *BibliographyRecord {
	return &BibliographyRecord{
		ID:            bib.ID.String(),
		BibIndex:      bib.BibIndex,
		Code:          bib.Code,
		Type:          bib.Type,
		Title:         bib.Title,
		Author:        bib.Author,
		Publisher:     bib.Publisher,
		ISBN:          bib.ISBN,
		PublishedDate: bib.PublishedDate.Format(time.RFC3339),
	}
}

// CSVBibliographyRepository implements domain.BibliographyRepository using a CSV file.
type CSVBibliographyRepository struct {
	FilePath string
}

func NewCSVBibliographyRepository(filePath string) *CSVBibliographyRepository {
	return &CSVBibliographyRepository{FilePath: filePath}
}

// Save implements domain.BibliographyRepository.Save
// Potential race condition: This method reads all records, modifies them, and writes them back
// without any locking mechanism. Acceptable for single-user CLI usage, but consider file locking
// or using a database with proper transaction support for production use.
func (r *CSVBibliographyRepository) Save(b *domain.Bibliography) error {
	records, err := ReadCSV(r.FilePath)
	if err != nil {
		return err
	}

	// Skip header
	if len(records) > 0 {
		records = records[1:]
	}

	iter := NewCSVRecordIterator(records, 0, 0)
	var all []*domain.Bibliography

	for iter.Next() {
		record := iter.Record()
		if len(record) < 9 {
			continue
		}
		bibRecord := &BibliographyRecord{
			ID:            record[0],
			BibIndex:      record[1],
			Code:          record[2],
			Type:          record[3],
			Title:         record[4],
			Author:        record[5],
			Publisher:     record[6],
			ISBN:          record[7],
			PublishedDate: record[8],
		}
		bib, err := recordToBibliography(bibRecord)
		if err != nil {
			slog.Error("Failed to convert bibliography record", "err", err)
			continue
		}
		all = append(all, bib)
	}

	if iter.Err() != nil {
		return iter.Err()
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

func (r *CSVBibliographyRepository) FindAll(limit, offset int) ([]*domain.Bibliography, error) {
	records, err := ReadCSV(r.FilePath)
	if err != nil {
		return nil, err
	}

	// Skip header
	if len(records) > 0 {
		records = records[1:]
	}

	iter := NewCSVRecordIterator(records, limit, offset)
	var bibliographies []*domain.Bibliography

	for iter.Next() {
		record := iter.Record()
		if len(record) < 9 {
			continue
		}

		bibRecord := &BibliographyRecord{
			ID:            record[0],
			BibIndex:      record[1],
			Code:          record[2],
			Type:          record[3],
			Title:         record[4],
			Author:        record[5],
			Publisher:     record[6],
			ISBN:          record[7],
			PublishedDate: record[8],
		}

		bib, err := recordToBibliography(bibRecord)
		if err != nil {
			slog.Error("Failed to convert bibliography record", "err", err)
			continue
		}

		bibliographies = append(bibliographies, bib)
	}

	return bibliographies, iter.Err()
}

// FindByBibIndex implements domain.BibliographyRepository.FindByBibIndex
func (r *CSVBibliographyRepository) FindByBibIndex(bibIndex string) (*domain.Bibliography, error) {
	records, err := ReadCSV(r.FilePath)
	if err != nil {
		return nil, err
	}

	// Skip header
	if len(records) > 0 {
		records = records[1:]
	}

	iter := NewCSVRecordIterator(records, 0, 0)

	for iter.Next() {
		record := iter.Record()
		if len(record) < 9 {
			continue
		}
		// Optimization: Check BibIndex (index 1) before full conversion
		if record[1] == bibIndex {
			bibRecord := &BibliographyRecord{
				ID:            record[0],
				BibIndex:      record[1],
				Code:          record[2],
				Type:          record[3],
				Title:         record[4],
				Author:        record[5],
				Publisher:     record[6],
				ISBN:          record[7],
				PublishedDate: record[8],
			}
			return recordToBibliography(bibRecord)
		}
	}

	return nil, iter.Err()
}

// FindByID implements domain.BibliographyRepository.FindByID
func (r *CSVBibliographyRepository) FindByID(id domain.BibliographyID) (*domain.Bibliography, error) {
	records, err := ReadCSV(r.FilePath)
	if err != nil {
		return nil, err
	}

	// Skip header
	if len(records) > 0 {
		records = records[1:]
	}

	iter := NewCSVRecordIterator(records, 0, 0)
	idStr := id.String()

	for iter.Next() {
		record := iter.Record()
		if len(record) < 9 {
			continue
		}
		// Optimization: Check ID (index 0) before full conversion
		if record[0] == idStr {
			bibRecord := &BibliographyRecord{
				ID:            record[0],
				BibIndex:      record[1],
				Code:          record[2],
				Type:          record[3],
				Title:         record[4],
				Author:        record[5],
				Publisher:     record[6],
				ISBN:          record[7],
				PublishedDate: record[8],
			}
			return recordToBibliography(bibRecord)
		}
	}

	return nil, iter.Err()
}

func (r *CSVBibliographyRepository) writeAll(bibliographies []*domain.Bibliography) error {
	var records [][]string
	// Header
	records = append(records, []string{"ID", "BibIndex", "Code", "Type", "Title", "Author", "Publisher", "ISBN", "PublishedDate"})

	for _, b := range bibliographies {
		rec := bibliographyToRecord(b)
		record := []string{
			rec.ID,
			rec.BibIndex,
			rec.Code,
			rec.Type,
			rec.Title,
			rec.Author,
			rec.Publisher,
			rec.ISBN,
			rec.PublishedDate,
		}
		records = append(records, record)
	}
	return WriteCSV(r.FilePath, records)
}
