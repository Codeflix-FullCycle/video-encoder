package utils

import "encoding/json"

func IsJson(s string) error {
	var js struct{}
	return json.Unmarshal([]byte(s), &js)
}
