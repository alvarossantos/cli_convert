package main

import (
	"errors"
	"fmt"
	"os"
)

func validateFileYaml(path string) error {
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

	return nil
}
