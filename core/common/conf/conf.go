package conf

import (
	"gopkg.in/yaml.v3"
	"os"
)

// Load load file data into an interface
func Load(file string, v any) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}
