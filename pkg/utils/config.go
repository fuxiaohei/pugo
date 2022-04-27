package utils

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

// LoadTOMLFile loads toml file.
func LoadTOMLFile(configFile string, v interface{}) error {
	if _, err := toml.DecodeFile(configFile, v); err != nil {
		return err
	}
	return nil
}

// WriteTOMLFile writes content to a TOML file
func WriteTOMLFile(path string, data interface{}) error {
	buffer := bytes.NewBuffer(nil)
	if err := toml.NewEncoder(buffer).Encode(data); err != nil {
		return err
	}
	return WriteFile(path, buffer.Bytes())
}
