package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	// Check if the Authorization header is present
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing") // No token provided
	}

	// Check if the header starts with "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", fmt.Errorf("invalid authorization header format")
	}

	// Extract the token
	token := authHeader[7:]
	token = strings.TrimSpace(token)
	if token == "" {
		return "", fmt.Errorf("token missing from authorization header")
	}
	return token, nil
}
