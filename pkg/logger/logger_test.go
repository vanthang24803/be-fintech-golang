package logger

import (
	"os"
	"testing"
)

func TestInitLoggerDevelopmentAndProduction(t *testing.T) {
	t.Parallel()

	prevEnv, hadEnv := os.LookupEnv("APP_ENV")
	defer func() {
		if hadEnv {
			_ = os.Setenv("APP_ENV", prevEnv)
		} else {
			_ = os.Unsetenv("APP_ENV")
		}
	}()

	_ = os.Unsetenv("APP_ENV")
	InitLogger()
	if Log == nil {
		t.Fatal("expected logger to be initialized in development mode")
	}

	if err := os.Setenv("APP_ENV", "production"); err != nil {
		t.Fatalf("Setenv(APP_ENV): %v", err)
	}
	InitLogger()
	if Log == nil {
		t.Fatal("expected logger to be initialized in production mode")
	}
}
