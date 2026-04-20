package service

import (
	"testing"

	"github.com/maynguyen24/sever/internal/models"
)

func TestFIDOUser_Methods(t *testing.T) {
	u := newFIDOUser(42, &models.Device{})
	if string(u.WebAuthnID()) != "42" {
		t.Fatalf("expected WebAuthnID '42', got '%s'", u.WebAuthnID())
	}
	if u.WebAuthnName() != "42" {
		t.Fatalf("expected WebAuthnName '42', got '%s'", u.WebAuthnName())
	}
	if u.WebAuthnDisplayName() != "42" {
		t.Fatalf("expected WebAuthnDisplayName '42', got '%s'", u.WebAuthnDisplayName())
	}
	if len(u.WebAuthnCredentials()) != 0 {
		t.Fatalf("expected empty credentials, got %d", len(u.WebAuthnCredentials()))
	}
}

func TestFIDOUser_WithCredentials(t *testing.T) {
	credID := "dGVzdA"    // base64url for "test"
	pubKey := "dGVzdA==" // base64 for "test"
	device := &models.Device{
		FIDOCredentialID: &credID,
		FIDOPublicKey:    &pubKey,
		FIDOSignCount:    5,
	}
	u := newFIDOUser(1, device)
	if len(u.WebAuthnCredentials()) != 1 {
		t.Fatalf("expected 1 credential, got %d", len(u.WebAuthnCredentials()))
	}
}

func TestFIDOUser_WithInvalidCredentials(t *testing.T) {
	// Invalid base64url should result in no credentials
	invalid := "!!invalid!!"
	device := &models.Device{
		FIDOCredentialID: &invalid,
		FIDOPublicKey:    &invalid,
	}
	u := newFIDOUser(1, device)
	if len(u.WebAuthnCredentials()) != 0 {
		t.Fatalf("expected 0 credentials for invalid base64, got %d", len(u.WebAuthnCredentials()))
	}
}
