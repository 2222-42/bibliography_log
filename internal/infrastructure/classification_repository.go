package infrastructure

import (
	"bibliography_log/internal/domain"
	"fmt"
	"log/slog"
	"strconv"
)

// ClassificationRecord represents a classification record for CSV persistence.
type ClassificationRecord struct {
	ID      string
	CodeNum string
	Name    string
}

// recordToClassification converts a ClassificationRecord to a domain.Classification.
func recordToClassification(rec *ClassificationRecord) (*domain.Classification, error) {
	id, err := domain.ParseClassificationID(rec.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse classification ID: %w", err)
	}

	codeNum, err := strconv.Atoi(rec.CodeNum)
	if err != nil {
		return nil, fmt.Errorf("failed to parse code number: %w", err)
	}

	return &domain.Classification{
		ID:      id,
		CodeNum: codeNum,
		Name:    rec.Name,
	}, nil
}

// classificationToRecord converts a domain.Classification to a ClassificationRecord.
func classificationToRecord(class *domain.Classification) *ClassificationRecord {
	return &ClassificationRecord{
		ID:      class.ID.String(),
		CodeNum: strconv.Itoa(class.CodeNum),
		Name:    class.Name,
	}
}

// CSVClassificationRepository implements domain.ClassificationRepository using a CSV file.
type CSVClassificationRepository struct {
	FilePath string
}

func NewCSVClassificationRepository(filePath string) *CSVClassificationRepository {
	return &CSVClassificationRepository{FilePath: filePath}
}

// Save implements domain.ClassificationRepository.Save
// Potential race condition: This method reads all records, modifies them, and writes them back
// without any locking mechanism. Acceptable for single-user CLI usage, but consider file locking
// or using a database with proper transaction support for production use.
func (r *CSVClassificationRepository) Save(c *domain.Classification) error {
	records, err := ReadCSV(r.FilePath)
	if err != nil {
		return err
	}

	// Skip header
	if len(records) > 0 {
		records = records[1:]
	}

	iter := NewCSVRecordIterator(records, 0, 0)
	var all []*domain.Classification

	for iter.Next() {
		record := iter.Record()
		if len(record) < 3 {
			continue
		}
		classRecord := &ClassificationRecord{
			ID:      record[0],
			CodeNum: record[1],
			Name:    record[2],
		}
		class, err := recordToClassification(classRecord)
		if err != nil {
			slog.Error("Failed to convert classification record", "err", err)
			continue
		}
		all = append(all, class)
	}

	if iter.Err() != nil {
		return iter.Err()
	}

	updated := false
	for i, existing := range all {
		if existing.ID == c.ID {
			all[i] = c
			updated = true
			break
		}
	}
	if !updated {
		all = append(all, c)
	}

	return r.writeAll(all)
}

func (r *CSVClassificationRepository) FindAll(limit, offset int) ([]*domain.Classification, error) {
	records, err := ReadCSV(r.FilePath)
	if err != nil {
		return nil, err
	}

	// Skip header
	if len(records) > 0 {
		records = records[1:]
	}

	iter := NewCSVRecordIterator(records, limit, offset)
	var classifications []*domain.Classification

	for iter.Next() {
		record := iter.Record()
		if len(record) < 3 {
			continue
		}

		classRecord := &ClassificationRecord{
			ID:      record[0],
			CodeNum: record[1],
			Name:    record[2],
		}

		class, err := recordToClassification(classRecord)
		if err != nil {
			slog.Error("Failed to convert classification record", "err", err)
			continue
		}

		classifications = append(classifications, class)
	}

	return classifications, iter.Err()
}

func (r *CSVClassificationRepository) FindByCodeNum(codeNum int) (*domain.Classification, error) {
	records, err := ReadCSV(r.FilePath)
	if err != nil {
		return nil, err
	}

	// Skip header
	if len(records) > 0 {
		records = records[1:]
	}

	iter := NewCSVRecordIterator(records, 0, 0)
	codeNumStr := strconv.Itoa(codeNum)

	for iter.Next() {
		record := iter.Record()
		if len(record) < 3 {
			continue
		}
		// Optimization: Check CodeNum (index 1) before full conversion
		if record[1] == codeNumStr {
			classRecord := &ClassificationRecord{
				ID:      record[0],
				CodeNum: record[1],
				Name:    record[2],
			}
			return recordToClassification(classRecord)
		}
	}

	return nil, iter.Err()
}

func (r *CSVClassificationRepository) writeAll(classifications []*domain.Classification) error {
	var records [][]string
	records = append(records, []string{"ID", "CodeNum", "Name"})

	for _, c := range classifications {
		rec := classificationToRecord(c)
		record := []string{
			rec.ID,
			rec.CodeNum,
			rec.Name,
		}
		records = append(records, record)
	}
	return WriteCSV(r.FilePath, records)
}
