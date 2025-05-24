package util

import "encoding/json"

// KeysFromMap extracts and returns a slice of keys from the given map.
func KeysFromMap(m map[string]string) []string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// ToJSON takes a map of string key-value pairs and returns a pretty-printed JSON string.
func ToJSON(config map[string]string) (string, error) {
	out, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
