package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT_and_ValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "mytestsecret"
	expiresIn := time.Minute

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	returnedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if returnedID != userID {
		t.Errorf("Expected %v, got %v", userID, returnedID)
	}
}
