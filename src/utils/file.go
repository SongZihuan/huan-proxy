package utils

import (
	"errors"
	"os"
)

func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}

	return s.IsDir()
}

func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !s.IsDir()
}

func MakeDir(path string) error {
	if IsExists(path) && IsDir(path) {
		return nil
	}

	return os.MkdirAll(path, os.ModePerm)
}
