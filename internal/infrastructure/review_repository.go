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
		id, _ := uuid.Parse(record[0])
		bookID, _ := uuid.Parse(record[1])
		createdAt, _ := time.Parse(time.RFC3339, record[4])
		updatedAt, _ := time.Parse(time.RFC3339, record[5])

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
