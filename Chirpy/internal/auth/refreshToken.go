package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	// Generate a random 32-byte token
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	// Encode the token as a hexadecimal string
	token := hex.EncodeToString(tokenBytes)
	return token, nil
}
