package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Get API key from the headers of HTTP request
// Example:
// Authorization: ApiKey {insert key}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")

	if val == "" {
		return "", errors.New("no authentication found")
	}

	vals := strings.Split(val, " ") // split value on space

	if len(vals) != 2 {
		return "", errors.New("malformed authentication header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("malformed authentication header")
	}

	return vals[1], nil
}
