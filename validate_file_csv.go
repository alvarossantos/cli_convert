package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
)

func validateFileCSV(path string, delimiter rune) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("file not exist")
		} else {
			return fmt.Errorf("access error : %v", err)
		}
	}

	if info.IsDir() {
		return errors.New("is a directory")
	}

	if info.Size() == 0 {
		return errors.New("file is empty")
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimiter
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}

	if len(records) == 0 {
		return errors.New("CSV file is empty")
	}

	expectedCols := len(records[0])

	if expectedCols == 1 {
		return fmt.Errorf("CSV appears to have only 1 column. Wrong delimiter? (current: '%c')", delimiter)
	}

	for i, record := range records[1:] {
		if len(record) != expectedCols {
			return fmt.Errorf("row %d: expected %d columns, but got %d", i+2, expectedCols, len(record))
		}

		for j, value := range record {
			if strings.TrimSpace(value) == "" {
				return fmt.Errorf("row %d: column '%s' is empty", i+2, records[0][j])
			}
		}
	}

	return nil
}
