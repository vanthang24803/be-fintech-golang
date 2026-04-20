package jwt

import (
	"testing"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/maynguyen24/sever/configs"
)

func testJWTConfig() *configs.Config {
	return &configs.Config{
		JWTSecret:        "test-jwt-secret-key",
		JWTRefreshSecret: "test-refresh-secret-key",
	}
}

func TestGenerateTokenPair(t *testing.T) {
	t.Parallel()

	cfg := testJWTConfig()
	accessToken, refreshToken, err := GenerateTokenPair(42, false, cfg)
	if err != nil {
		t.Fatalf("GenerateTokenPair returned error: %v", err)
	}
	if accessToken == "" || refreshToken == "" {
		t.Fatal("expected non-empty token pair")
	}

	// Verify access token claims
	claims := &TokenClaims{}
	token, err := gojwt.ParseWithClaims(accessToken, claims, func(token *gojwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("access token invalid: %v", err)
	}
	if claims.UserID != 42 {
		t.Fatalf("expected user ID 42, got %d", claims.UserID)
	}
	if claims.FIDOVerified {
		t.Fatal("expected FIDOVerified false")
	}

	// Verify refresh token claims
	refreshClaims := &TokenClaims{}
	rToken, err := gojwt.ParseWithClaims(refreshToken, refreshClaims, func(token *gojwt.Token) (interface{}, error) {
		return []byte(cfg.JWTRefreshSecret), nil
	})
	if err != nil || !rToken.Valid {
		t.Fatalf("refresh token invalid: %v", err)
	}
	if refreshClaims.UserID != 42 {
		t.Fatalf("expected user ID 42 in refresh token, got %d", refreshClaims.UserID)
	}
}

func TestGenerateTokenPair_FIDOVerified(t *testing.T) {
	t.Parallel()

	cfg := testJWTConfig()
	accessToken, _, err := GenerateTokenPair(99, true, cfg)
	if err != nil {
		t.Fatalf("GenerateTokenPair returned error: %v", err)
	}

	claims := &TokenClaims{}
	_, err = gojwt.ParseWithClaims(accessToken, claims, func(token *gojwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		t.Fatalf("access token invalid: %v", err)
	}
	if !claims.FIDOVerified {
		t.Fatal("expected FIDOVerified true")
	}
	if claims.UserID != 99 {
		t.Fatalf("expected user ID 99, got %d", claims.UserID)
	}
}

func TestGenerateTokenPair_DifferentTokens(t *testing.T) {
	t.Parallel()

	cfg := testJWTConfig()
	at1, rt1, err := GenerateTokenPair(1, false, cfg)
	if err != nil {
		t.Fatalf("first pair error: %v", err)
	}
	at2, rt2, err := GenerateTokenPair(2, false, cfg)
	if err != nil {
		t.Fatalf("second pair error: %v", err)
	}

	if at1 == at2 {
		t.Fatal("expected different access tokens for different users")
	}
	if rt1 == rt2 {
		t.Fatal("expected different refresh tokens for different users")
	}
}

func TestGenerateAccessTokenFIDO(t *testing.T) {
	t.Parallel()

	cfg := testJWTConfig()
	accessToken, err := GenerateAccessTokenFIDO(77, cfg)
	if err != nil {
		t.Fatalf("GenerateAccessTokenFIDO returned error: %v", err)
	}
	if accessToken == "" {
		t.Fatal("expected non-empty access token")
	}

	claims := &TokenClaims{}
	_, err = gojwt.ParseWithClaims(accessToken, claims, func(token *gojwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		t.Fatalf("FIDO access token invalid: %v", err)
	}
	if claims.UserID != 77 {
		t.Fatalf("expected user ID 77, got %d", claims.UserID)
	}
	if !claims.FIDOVerified {
		t.Fatal("expected FIDOVerified true")
	}
}

func TestGenerateAccessTokenFIDO_IsSignedWithAccessSecret(t *testing.T) {
	t.Parallel()

	cfg := testJWTConfig()
	accessToken, err := GenerateAccessTokenFIDO(10, cfg)
	if err != nil {
		t.Fatalf("GenerateAccessTokenFIDO returned error: %v", err)
	}

	// Must be verifiable with JWTSecret (not refresh secret)
	claims := &TokenClaims{}
	_, err = gojwt.ParseWithClaims(accessToken, claims, func(token *gojwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		t.Fatalf("FIDO token should be verifiable with JWTSecret: %v", err)
	}

	// Must NOT be verifiable with refresh secret
	_, err = gojwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *gojwt.Token) (interface{}, error) {
		return []byte(cfg.JWTRefreshSecret), nil
	})
	if err == nil {
		t.Fatal("FIDO token should not be verifiable with JWTRefreshSecret")
	}
}
