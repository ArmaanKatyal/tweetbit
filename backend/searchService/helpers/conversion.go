package helpers

import "encoding/base64"

// convert string to bool
func StringToBool(s string) bool {
	return s == "true"
}

// base64 decode a string
func Base64Decode(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
