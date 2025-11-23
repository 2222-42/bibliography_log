package infrastructure

import (
	"bibliography_log/internal/domain"
	"fmt"
	"log/slog"
	"time"
)

// ReviewRecord represents a review record for CSV persistence.
type ReviewRecord struct {
	ID        string
	BookID    string
	Goals     string
	Summary   string
	CreatedAt string
	UpdatedAt string
}

// recordToReview converts a ReviewRecord to a domain.Review.
func recordToReview(rec *ReviewRecord) (*domain.Review, error) {
	id, err := domain.ParseReviewID(rec.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse review ID: %w", err)
	}

	bookID, err := domain.ParseBibliographyID(rec.BookID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse book ID: %w", err)
	}

	createdAt, err := time.Parse(time.RFC3339, rec.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created at: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, rec.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated at: %w", err)
	}

	return &domain.Review{
		ID:        id,
		BookID:    bookID,
		Goals:     rec.Goals,
		Summary:   rec.Summary,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// reviewToRecord converts a domain.Review to a ReviewRecord.
func reviewToRecord(rev *domain.Review) *ReviewRecord {
	return &ReviewRecord{
		ID:        rev.ID.String(),
		BookID:    rev.BookID.String(),
		Goals:     rev.Goals,
		Summary:   rev.Summary,
		CreatedAt: rev.CreatedAt.Format(time.RFC3339),
		UpdatedAt: rev.UpdatedAt.Format(time.RFC3339),
	}
}

// CSVReviewRepository implements domain.ReviewRepository using a CSV file.
type CSVReviewRepository struct {
	FilePath string
}

func NewCSVReviewRepository(filePath string) *CSVReviewRepository {
	return &CSVReviewRepository{FilePath: filePath}
}

// Save implements domain.ReviewRepository.Save
// This contains potential race condition. But, it is not a problem in this cli application.
func (r *CSVReviewRepository) Save(review *domain.Review) error {
	records, err := ReadCSV(r.FilePath)
	if err != nil {
		return err
	}

	// Skip header
	if len(records) > 0 {
		records = records[1:]
	}

	iter := NewCSVRecordIterator(records, 0, 0)
	var all []*domain.Review

	for iter.Next() {
		record := iter.Record()
		if len(record) < 6 {
			continue
		}
		revRecord := &ReviewRecord{
			ID:        record[0],
			BookID:    record[1],
			Goals:     record[2],
			Summary:   record[3],
			CreatedAt: record[4],
			UpdatedAt: record[5],
		}
		rev, err := recordToReview(revRecord)
		if err != nil {
			slog.Error("Failed to convert review record", "err", err)
			continue
		}
		all = append(all, rev)
	}

	if iter.Err() != nil {
		return iter.Err()
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

func (r *CSVReviewRepository) FindAll(limit, offset int) ([]*domain.Review, error) {
	records, err := ReadCSV(r.FilePath)
	if err != nil {
		return nil, err
	}

	// Skip header
	if len(records) > 0 {
		records = records[1:]
	}

	iter := NewCSVRecordIterator(records, limit, offset)
	var reviews []*domain.Review

	for iter.Next() {
		record := iter.Record()
		if len(record) < 6 {
			continue
		}

		revRecord := &ReviewRecord{
			ID:        record[0],
			BookID:    record[1],
			Goals:     record[2],
			Summary:   record[3],
			CreatedAt: record[4],
			UpdatedAt: record[5],
		}

		rev, err := recordToReview(revRecord)
		if err != nil {
			slog.Error("Failed to convert review record", "err", err)
			continue
		}

		reviews = append(reviews, rev)
	}

	return reviews, iter.Err()
}

// FindByID implements domain.ReviewRepository.FindByID
func (r *CSVReviewRepository) FindByID(id domain.ReviewID) (*domain.Review, error) {
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
		if len(record) < 6 {
			continue
		}
		// Optimization: Check ID (index 0) before full conversion
		if record[0] == idStr {
			revRecord := &ReviewRecord{
				ID:        record[0],
				BookID:    record[1],
				Goals:     record[2],
				Summary:   record[3],
				CreatedAt: record[4],
				UpdatedAt: record[5],
			}
			return recordToReview(revRecord)
		}
	}

	return nil, iter.Err()
}

func (r *CSVReviewRepository) FindByBookID(bookID domain.BibliographyID) ([]*domain.Review, error) {
	records, err := ReadCSV(r.FilePath)
	if err != nil {
		return nil, err
	}

	// Skip header
	if len(records) > 0 {
		records = records[1:]
	}

	iter := NewCSVRecordIterator(records, 0, 0)
	bookIDStr := bookID.String()
	var matches []*domain.Review

	for iter.Next() {
		record := iter.Record()
		if len(record) < 6 {
			continue
		}
		// Optimization: Check BookID (index 1) before full conversion
		if record[1] == bookIDStr {
			revRecord := &ReviewRecord{
				ID:        record[0],
				BookID:    record[1],
				Goals:     record[2],
				Summary:   record[3],
				CreatedAt: record[4],
				UpdatedAt: record[5],
			}
			rev, err := recordToReview(revRecord)
			if err != nil {
				slog.Error("Failed to convert review record", "err", err)
				continue
			}
			matches = append(matches, rev)
		}
	}

	return matches, iter.Err()
}

func (r *CSVReviewRepository) writeAll(reviews []*domain.Review) error {
	var records [][]string
	records = append(records, []string{"ID", "BookID", "Goals", "Summary", "CreatedAt", "UpdatedAt"})

	for _, review := range reviews {
		rec := reviewToRecord(review)
		record := []string{
			rec.ID,
			rec.BookID,
			rec.Goals,
			rec.Summary,
			rec.CreatedAt,
			rec.UpdatedAt,
		}
		records = append(records, record)
	}
	return WriteCSV(r.FilePath, records)
}
