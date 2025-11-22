package infrastructure

import (
	"strconv"

	"bibliography_log/internal/domain"

	"github.com/google/uuid"
)

// CSVBibClassificationRepository implements domain.BibClassificationRepository using a CSV file.
type CSVBibClassificationRepository struct {
	FilePath string
}

func NewCSVBibClassificationRepository(filePath string) *CSVBibClassificationRepository {
	return &CSVBibClassificationRepository{FilePath: filePath}
}

// Save implements domain.BibClassificationRepository.Save
// Potential race condition: This method reads all records, modifies them, and writes them back
// without any locking mechanism. Acceptable for single-user CLI usage, but consider file locking
// or using a database with proper transaction support for production use.
func (r *CSVBibClassificationRepository) Save(c *domain.BibClassification) error {
	all, err := r.FindAll()
	if err != nil {
		return err
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

func (r *CSVBibClassificationRepository) FindAll() ([]*domain.BibClassification, error) {
	records, err := ReadCSV(r.FilePath)
	if err != nil {
		return nil, err
	}

	var classifications []*domain.BibClassification
	if len(records) > 0 {
		records = records[1:]
	}

	for _, record := range records {
		if len(record) < 3 {
			continue
		}
		id, err := uuid.Parse(record[0])
		if err != nil {
			continue // Skip records with invalid UUIDs
		}
		codeNum, err := strconv.Atoi(record[1])
		if err != nil {
			continue // Skip records with invalid code numbers
		}

		classifications = append(classifications, &domain.BibClassification{
			ID:      id,
			CodeNum: codeNum,
			Name:    record[2],
		})
	}
	return classifications, nil
}

func (r *CSVBibClassificationRepository) FindByCodeNum(codeNum int) (*domain.BibClassification, error) {
	all, err := r.FindAll()
	if err != nil {
		return nil, err
	}
	for _, c := range all {
		if c.CodeNum == codeNum {
			return c, nil
		}
	}
	return nil, nil
}

func (r *CSVBibClassificationRepository) writeAll(classifications []*domain.BibClassification) error {
	var records [][]string
	records = append(records, []string{"ID", "CodeNum", "Name"})

	for _, c := range classifications {
		record := []string{
			c.ID.String(),
			strconv.Itoa(c.CodeNum),
			c.Name,
		}
		records = append(records, record)
	}
	return WriteCSV(r.FilePath, records)
}
