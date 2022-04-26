package utils

import "github.com/BurntSushi/toml"

// LoadTOML loads toml file.
func LoadTOML(configFile string, v interface{}) error {
	if _, err := toml.DecodeFile(configFile, v); err != nil {
		return err
	}
	return nil
}
