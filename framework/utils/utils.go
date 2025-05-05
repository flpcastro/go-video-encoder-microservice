package utils

import "encoding/json"

func IsJSON(s string) error {
	var js struct{}

	err := json.Unmarshal([]byte(s), &js)
	if err != nil {
		return err
	}

	return nil
}
