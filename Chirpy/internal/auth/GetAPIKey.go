package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers["Authorization"]
	if len(authHeader) == 0 {
		return "", fmt.Errorf("authorization header not found")
	}
	if len(authHeader) > 1 {
		return "", fmt.Errorf("multiple authorization headers found")
	}
	authValue := authHeader[0]
	if !strings.HasPrefix(authValue, "ApiKey ") {
		return "", fmt.Errorf("authorization header does not start with 'ApiKey '")
	}
	apiKey := strings.TrimPrefix(authValue, "ApiKey ")
	if apiKey == "" {
		return "", fmt.Errorf("API key is empty")
	}
	return apiKey, nil
}
