package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
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
	reader.LazyQuotes = true
	expectedCols := 0
	for {
		read, err := reader.Read()
		if err == io.EOF {
			break
		}
		
		if err != nil {
			return err
		}

		if expectedCols == 0 {
			expectedCols = len(read)
		} else {
			if len(read) != expectedCols {
				return fmt.Errorf("inconsistent number of columns: expected %d, got %d", expectedCols, len(read))
			}
		}
	}

	return nil
}
