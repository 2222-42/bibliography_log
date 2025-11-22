package infrastructure

import (
	"encoding/csv"
	"log"
	"os"
)

// ReadCSV reads all records from a CSV file.
// It returns the records including the header if present (caller handles header).
func ReadCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if os.IsNotExist(err) {
		return [][]string{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}

// WriteCSV writes records to a CSV file.
// It overwrites the file.
func WriteCSV(filePath string, records [][]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}()

	writer := csv.NewWriter(file)

	if err := writer.WriteAll(records); err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}
