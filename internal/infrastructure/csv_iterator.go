package infrastructure

// CSVRecordIterator provides an iterator interface for CSV records with pagination support.
type CSVRecordIterator struct {
	records [][]string
	current int
	limit   int
	offset  int
	err     error
	start   int // Starting position after offset
	end     int // Ending position (exclusive)
}

// NewCSVRecordIterator creates a new iterator from CSV records.
// The records slice should already have the header removed if applicable.
// limit <= 0 means no limit (return all records after offset).
// offset is the number of records to skip from the beginning.
func NewCSVRecordIterator(records [][]string, limit, offset int) *CSVRecordIterator {
	start := offset
	if start < 0 {
		start = 0
	}
	if start > len(records) {
		start = len(records)
	}

	end := len(records)
	if limit > 0 {
		end = start + limit
		if end > len(records) {
			end = len(records)
		}
	}

	return &CSVRecordIterator{
		records: records,
		current: start - 1, // Start before the first record
		limit:   limit,
		offset:  offset,
		start:   start,
		end:     end,
	}
}

// Next advances to the next record and returns true if a record is available.
func (it *CSVRecordIterator) Next() bool {
	it.current++
	return it.current < it.end
}

// Record returns the current record.
// Should only be called after Next() returns true.
func (it *CSVRecordIterator) Record() []string {
	if it.current < 0 || it.current >= len(it.records) {
		return nil
	}
	return it.records[it.current]
}

// Err returns any error that occurred during iteration.
// Currently, CSV iteration doesn't produce errors, but this method
// is provided for consistency with other iterator patterns.
func (it *CSVRecordIterator) Err() error {
	return it.err
}

// Reset resets the iterator to the beginning (respecting offset).
func (it *CSVRecordIterator) Reset() {
	it.current = it.start - 1
	it.err = nil
}

// Count returns the total number of records available in this iteration
// (after applying offset and limit).
func (it *CSVRecordIterator) Count() int {
	return it.end - it.start
}
