package utils

import "strings"

// GetKeyValue splits a callback query data into key and value.
// The data is expected to have the format "key=value".
func GetKeyValue(data string) (string, string) {
	parts := strings.SplitN(data, "=", 2) // Split the string at the first "="
	if len(parts) != 2 {
		return "", "" // Return empty strings if the format is invalid
	}
	return parts[0], parts[1] // Return key and value
}
