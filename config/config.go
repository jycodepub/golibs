package config

import (
	"encoding/json"
	"os"
)

func Load(filepath string, values interface{}) error {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, values)
}
