package utils

import (
	"bytes"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
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

// WriteYAMLFile writes content to a yaml file
func WriteYAMLFile(path string, data interface{}) error {
	dataBytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return WriteFile(path, dataBytes)
}

// LoadYAMLFile loads yaml file.
func LoadYAMLFile(configFile string, v interface{}) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, v)
}
