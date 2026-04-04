package configs

import "os"

// Config holds all centralized application configurations
type Config struct {
	AppEnv           string
	Port             string
	DatabaseURL      string
	JWTSecret        string
	JWTRefreshSecret string
	GoogleClientID   string
	GoogleClientSecret string
	GoogleRedirectURL string
	FirebaseCredentialsPath string
	RedisAddr               string
	RedisPassword           string
}

// LoadConfig initializes the configuration struct from environment variables
func LoadConfig() *Config {
	return &Config{
		AppEnv:           getEnv("APP_ENV", "development"),
		Port:             getEnv("PORT", "8386"), // Default to 8386 per user setup
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		JWTSecret:        getEnv("JWT_SECRET", "super_secret_jwt_key_for_dev_only"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "super_secret_refresh_key_for_dev_only"),
		GoogleClientID:   os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL: getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8386/api/v1/auth/google/callback"),
		FirebaseCredentialsPath: getEnv("FIREBASE_SERVICE_ACCOUNT_JSON", "configs/firebase-service-account.json"),
		RedisAddr:               getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:           getEnv("REDIS_PASSWORD", ""),
	}
}

// getEnv handles fallback logic for environment variables
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
