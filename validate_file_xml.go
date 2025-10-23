package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
)

func validateFileXML(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not exist")
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

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	var v interface{}
	if err := xml.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("invalid XML format: %v", err)
	}
	return nil
}
