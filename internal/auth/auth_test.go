package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	// Setup
	secret := "your-test-secret"
	userID := uuid.New()

	// Test cases
	t.Run("valid token", func(t *testing.T) {
		// Create a token that should work
		token, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("failed to create token: %v", err)
		}

		// Validate it
		gotUID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if gotUID != userID {
			t.Errorf("expected user ID %v, got %v", userID, gotUID)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		// Create a token that's already expired
		token, _ := MakeJWT(userID, secret, -time.Hour)

		_, err := ValidateJWT(token, secret)
		if err == nil {
			t.Error("expected error for expired token, got nil")
		}
	})

	t.Run("wrong secret", func(t *testing.T) {
		// Create token with one secret
		token, _ := MakeJWT(userID, secret, time.Hour)

		// Try to validate with different secret
		_, err := ValidateJWT(token, "wrong-secret")
		if err == nil {
			t.Error("expected error for wrong secret, got nil")
		}
	})
}
