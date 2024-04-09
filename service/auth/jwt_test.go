package auth

import (
	"testing"

	"github.com/google/uuid"
)

func TestCreateJWT(t *testing.T) {
	secret := []byte("secret")

	token, err := CreateJWT(secret, uuid.New().String())
	if err != nil {
		t.Errorf("error creating JWT: %v", err)
	}

	if token == "" {
		t.Error("expected token to be not empty")
	}
}
