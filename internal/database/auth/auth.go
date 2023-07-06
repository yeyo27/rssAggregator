package auth

import (
	"errors"
	"net/http"
	"strings"
)

// this function extracts the api key from the headers of the HTTP request
// Example: Authorization: "ApiKey {apikey value}"
func GetAPIKey(headers http.Header) (string, error) {
	value := headers.Get("Authorization")
	if value == "" {
		return "", errors.New("No authentication info found")
	}

	values := strings.Split(value, " ")
	if len(values) != 2 {
		return "", errors.New("Incorrect authentication header")
	}

	if values[0] != "ApiKey" {
		return "", errors.New("Incorrect API Key format")
	}

	return values[1], nil
}