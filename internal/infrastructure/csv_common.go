package infrastructure

import (
	"encoding/csv"
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
	defer file.Close()

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
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.WriteAll(records)
}
