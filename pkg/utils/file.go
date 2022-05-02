package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

// CopyFile copies a file
func CopyFile(src, dst string) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// create subdir if necessary
	dir := filepath.Dir(dst)
	if !IsDirExist(dir) {
		if err := MkdirAll(dir); err != nil {
			return err
		}
	}

	// if exist, remove it
	if IsFileExist(dst) {
		if err := os.Remove(dst); err != nil {
			return err
		}
	}

	w, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer w.Close()
	w.ReadFrom(r)

	return nil
}

// IsTempFile checks if a file is a temporary file
func IsTempFile(fpath string) bool {
	ext := filepath.Ext(fpath)
	baseName := filepath.Base(fpath)
	istemp := strings.HasSuffix(ext, "~") ||
		(ext == ".swp") || // vim
		(ext == ".swx") || // vim
		(ext == ".tmp") || // generic temp file
		(ext == ".DS_Store") || // OSX Thumbnail
		baseName == "4913" || // vim
		strings.HasPrefix(ext, ".goutputstream") || // gnome
		strings.HasSuffix(ext, "jb_old___") || // intelliJ
		strings.HasSuffix(ext, "jb_tmp___") || // intelliJ
		strings.HasSuffix(ext, "jb_bak___") || // intelliJ
		strings.HasPrefix(ext, ".sb-") || // byword
		strings.HasPrefix(baseName, ".#") || // emacs
		strings.HasPrefix(baseName, "#") // emacs
	return istemp
}
