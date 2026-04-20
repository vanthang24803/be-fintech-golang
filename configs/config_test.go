package configs

import "testing"

func TestLoadConfigAndGetEnv(t *testing.T) {
	t.Setenv("PORT", "9999")
	t.Setenv("JWT_SECRET", "jwt-secret")
	t.Setenv("JWT_REFRESH_SECRET", "refresh-secret")
	t.Setenv("GOOGLE_REDIRECT_URL", "http://localhost/callback")
	t.Setenv("MINIO_USE_SSL", "true")

	cfg := LoadConfig()
	if cfg.Port != "9999" || cfg.JWTSecret != "jwt-secret" || cfg.JWTRefreshSecret != "refresh-secret" {
		t.Fatalf("unexpected config values: %+v", cfg)
	}
	if !cfg.MinioUseSSL {
		t.Fatal("expected MINIO_USE_SSL=true to be parsed")
	}
	if got := getEnv("NON_EXISTENT_ENV", "fallback"); got != "fallback" {
		t.Fatalf("expected fallback, got %q", got)
	}
}
