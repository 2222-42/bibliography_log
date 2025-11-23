package infrastructure

import (
	"testing"
)

func TestCSVRecordIterator_BasicIteration(t *testing.T) {
	records := [][]string{
		{"1", "a", "b"},
		{"2", "c", "d"},
		{"3", "e", "f"},
	}

	iter := NewCSVRecordIterator(records, 0, 0)

	// Test Count
	if iter.Count() != 3 {
		t.Errorf("Expected count 3, got %d", iter.Count())
	}

	// Test iteration
	count := 0
	for iter.Next() {
		record := iter.Record()
		if record == nil {
			t.Error("Expected non-nil record")
		}
		count++
	}

	if count != 3 {
		t.Errorf("Expected to iterate 3 times, got %d", count)
	}

	// Test Err
	if iter.Err() != nil {
		t.Errorf("Expected no error, got %v", iter.Err())
	}
}

func TestCSVRecordIterator_WithLimit(t *testing.T) {
	records := [][]string{
		{"1", "a"},
		{"2", "b"},
		{"3", "c"},
		{"4", "d"},
		{"5", "e"},
	}

	iter := NewCSVRecordIterator(records, 3, 0)

	if iter.Count() != 3 {
		t.Errorf("Expected count 3, got %d", iter.Count())
	}

	count := 0
	for iter.Next() {
		count++
	}

	if count != 3 {
		t.Errorf("Expected to iterate 3 times, got %d", count)
	}
}

func TestCSVRecordIterator_WithOffset(t *testing.T) {
	records := [][]string{
		{"1", "a"},
		{"2", "b"},
		{"3", "c"},
		{"4", "d"},
		{"5", "e"},
	}

	iter := NewCSVRecordIterator(records, 0, 2)

	if iter.Count() != 3 {
		t.Errorf("Expected count 3, got %d", iter.Count())
	}

	// First record should be the third one (index 2)
	if !iter.Next() {
		t.Fatal("Expected Next() to return true")
	}
	record := iter.Record()
	if record[0] != "3" {
		t.Errorf("Expected first record to be '3', got '%s'", record[0])
	}
}

func TestCSVRecordIterator_WithLimitAndOffset(t *testing.T) {
	records := [][]string{
		{"1", "a"},
		{"2", "b"},
		{"3", "c"},
		{"4", "d"},
		{"5", "e"},
	}

	iter := NewCSVRecordIterator(records, 2, 1)

	if iter.Count() != 2 {
		t.Errorf("Expected count 2, got %d", iter.Count())
	}

	// Should get records at index 1 and 2
	if !iter.Next() {
		t.Fatal("Expected Next() to return true")
	}
	record := iter.Record()
	if record[0] != "2" {
		t.Errorf("Expected first record to be '2', got '%s'", record[0])
	}

	if !iter.Next() {
		t.Fatal("Expected Next() to return true")
	}
	record = iter.Record()
	if record[0] != "3" {
		t.Errorf("Expected second record to be '3', got '%s'", record[0])
	}

	if iter.Next() {
		t.Error("Expected Next() to return false")
	}
}

func TestCSVRecordIterator_EmptyRecords(t *testing.T) {
	records := [][]string{}

	iter := NewCSVRecordIterator(records, 0, 0)

	if iter.Count() != 0 {
		t.Errorf("Expected count 0, got %d", iter.Count())
	}

	if iter.Next() {
		t.Error("Expected Next() to return false for empty records")
	}
}

func TestCSVRecordIterator_OffsetBeyondLength(t *testing.T) {
	records := [][]string{
		{"1", "a"},
		{"2", "b"},
	}

	iter := NewCSVRecordIterator(records, 0, 10)

	if iter.Count() != 0 {
		t.Errorf("Expected count 0, got %d", iter.Count())
	}

	if iter.Next() {
		t.Error("Expected Next() to return false")
	}
}

func TestCSVRecordIterator_Reset(t *testing.T) {
	records := [][]string{
		{"1", "a"},
		{"2", "b"},
		{"3", "c"},
	}

	iter := NewCSVRecordIterator(records, 0, 0)

	// Iterate through all records
	count := 0
	for iter.Next() {
		count++
	}

	if count != 3 {
		t.Errorf("Expected to iterate 3 times, got %d", count)
	}

	// Reset and iterate again
	iter.Reset()
	count = 0
	for iter.Next() {
		count++
	}

	if count != 3 {
		t.Errorf("Expected to iterate 3 times after reset, got %d", count)
	}
}

func TestCSVRecordIterator_ResetWithPagination(t *testing.T) {
	records := [][]string{
		{"1", "a"},
		{"2", "b"},
		{"3", "c"},
		{"4", "d"},
		{"5", "e"},
	}

	iter := NewCSVRecordIterator(records, 2, 1)

	// Iterate through paginated records
	count := 0
	for iter.Next() {
		count++
	}

	if count != 2 {
		t.Errorf("Expected to iterate 2 times, got %d", count)
	}

	// Reset and iterate again
	iter.Reset()
	count = 0
	var firstRecord string
	for iter.Next() {
		if count == 0 {
			firstRecord = iter.Record()[0]
		}
		count++
	}

	if count != 2 {
		t.Errorf("Expected to iterate 2 times after reset, got %d", count)
	}

	if firstRecord != "2" {
		t.Errorf("Expected first record after reset to be '2', got '%s'", firstRecord)
	}
}

func TestCSVRecordIterator_NegativeOffset(t *testing.T) {
	records := [][]string{
		{"1", "a"},
		{"2", "b"},
		{"3", "c"},
	}

	iter := NewCSVRecordIterator(records, 0, -5)

	// Negative offset should be treated as 0
	if iter.Count() != 3 {
		t.Errorf("Expected count 3, got %d", iter.Count())
	}

	if !iter.Next() {
		t.Fatal("Expected Next() to return true")
	}
	record := iter.Record()
	if record[0] != "1" {
		t.Errorf("Expected first record to be '1', got '%s'", record[0])
	}
}
