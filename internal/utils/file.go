package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// IsFileExist checks if a file exists
func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// IsDirExist checks if a directory exists
func IsDirExist(dir string) bool {
	fi, err := os.Stat(dir)
	if err != nil {
		return os.IsExist(err)
	}
	return fi.IsDir()
}

// MkdirAll creates a directory recursively
func MkdirAll(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// WriteFile writes content to a file
func WriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if !IsDirExist(dir) {
		if err := MkdirAll(dir); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(path, data, os.ModePerm)
}
