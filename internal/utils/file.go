package utils

import "os"

// IsFileExist checks if a file exists
func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// MkdirAll creates a directory recursively
func MkdirAll(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
