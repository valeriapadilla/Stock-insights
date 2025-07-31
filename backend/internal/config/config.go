package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Port           string
	Environment    string
	LogLevel       string
	DatabaseURL    string
	ExternalAPIURL string
	ExternalAPIKey string
	CacheTTL       time.Duration
	RateLimit      int
	AdminAPIKey    string
}

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	config := &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),

		DatabaseURL: getEnv("DATABASE_URL", ""),

		ExternalAPIURL: getEnv("EXTERNAL_API_URL", "https://api.karenai.click"),
		ExternalAPIKey: getEnv("EXTERNAL_API_KEY", ""),

		CacheTTL:  getEnvAsDuration("CACHE_TTL", 5*time.Minute),
		RateLimit: getEnvAsInt("RATE_LIMIT", 100),

		AdminAPIKey: getEnv("ADMIN_API_KEY", ""),
	}

	return config
}

func (c *Config) Validate() error {
	if c.Environment == "test" {
		return nil
	}

	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if c.Environment != "development" && c.Environment != "production" && c.Environment != "test" {
		return fmt.Errorf("ENVIRONMENT must be one of: development, production, test")
	}

	if c.RateLimit <= 0 {
		return fmt.Errorf("RATE_LIMIT must be greater than 0")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
