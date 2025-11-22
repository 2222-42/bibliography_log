package infrastructure

import (
	"bibliography_log/internal/domain"
	"time"

	"github.com/google/uuid"
)

// CSVReviewRepository implements domain.ReviewRepository using a CSV file.
type CSVReviewRepository struct {
	FilePath string
}

func NewCSVReviewRepository(filePath string) *CSVReviewRepository {
	return &CSVReviewRepository{FilePath: filePath}
}

// Save implements domain.ReviewRepository.Save
// This contains potential race conditoin. But, it is not a problem in this cli application.
func (r *CSVReviewRepository) Save(review *domain.Review) error {
	all, err := r.FindAll()
	if err != nil {
		return err
	}

	updated := false
	for i, existing := range all {
		if existing.ID == review.ID {
			all[i] = review
			updated = true
			break
		}
	}
	if !updated {
		all = append(all, review)
	}

	return r.writeAll(all)
}

func (r *CSVReviewRepository) FindAll() ([]*domain.Review, error) {
	records, err := ReadCSV(r.FilePath)
	if err != nil {
		return nil, err
	}

	var reviews []*domain.Review
	if len(records) > 0 {
		records = records[1:]
	}

	for _, record := range records {
		if len(record) < 6 {
			continue
		}
		id, err := uuid.Parse(record[0])
		if err != nil {
			continue // Skip records with invalid UUIDs
		}
		bookID, err := uuid.Parse(record[1])
		if err != nil {
			continue // Skip records with invalid book IDs
		}
		createdAt, err := time.Parse(time.RFC3339, record[4])
		if err != nil {
			continue // Skip records with invalid timestamps
		}
		updatedAt, err := time.Parse(time.RFC3339, record[5])
		if err != nil {
			continue // Skip records with invalid timestamps
		}

		reviews = append(reviews, &domain.Review{
			ID:        id,
			BookID:    bookID,
			Goals:     record[2],
			Summary:   record[3],
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}
	return reviews, nil
}

// FindByID implements domain.ReviewRepository.FindByID
// Performance Note: This method calls FindAll() which reads and parses the entire CSV file.
// For large datasets, consider implementing caching or using a database for production use.
func (r *CSVReviewRepository) FindByID(id uuid.UUID) (*domain.Review, error) {
	all, err := r.FindAll()
	if err != nil {
		return nil, err
	}
	for _, review := range all {
		if review.ID == id {
			return review, nil
		}
	}
	return nil, nil
}

func (r *CSVReviewRepository) FindByBookID(bookID uuid.UUID) ([]*domain.Review, error) {
	all, err := r.FindAll()
	if err != nil {
		return nil, err
	}
	var matches []*domain.Review
	for _, review := range all {
		if review.BookID == bookID {
			matches = append(matches, review)
		}
	}
	return matches, nil
}

func (r *CSVReviewRepository) writeAll(reviews []*domain.Review) error {
	var records [][]string
	records = append(records, []string{"ID", "BookID", "Goals", "Summary", "CreatedAt", "UpdatedAt"})

	for _, review := range reviews {
		record := []string{
			review.ID.String(),
			review.BookID.String(),
			review.Goals,
			review.Summary,
			review.CreatedAt.Format(time.RFC3339),
			review.UpdatedAt.Format(time.RFC3339),
		}
		records = append(records, record)
	}
	return WriteCSV(r.FilePath, records)
}
