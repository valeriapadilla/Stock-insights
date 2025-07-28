package config

import(
	"os"
	"strconv"
	"time"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct{
	Port string
	Environment string
	LogLevel string

	DatabaseURL string

	ExternalAPIURL string
	ExternalAPIKey string

	CacheTTL time.Duration

	RateLimit int

	AdminAPIKey string

	MetricsPort string
	EnableTracing bool
}

func Load() (*Config){
	if err := godotenv.Load(); err != nil{
		logrus.Warn("Error loading .env file: %v", err)
	}
	config := &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		
		DatabaseURL: getEnv("DATABASE_URL", ""),
		
		ExternalAPIURL: getEnv("EXTERNAL_API_URL", ""),
		ExternalAPIKey: getEnv("EXTERNAL_API_KEY", ""),
		
		CacheTTL: getEnvAsDuration("CACHE_TTL", 5*time.Minute),
		RateLimit: getEnvAsInt("RATE_LIMIT", 100),
		
		AdminAPIKey: getEnv("ADMIN_API_KEY", ""),
		
		MetricsPort: getEnv("METRICS_PORT", "9090"),
		EnableTracing: getEnvAsBool("ENABLE_TRACING", false),
	}
	return config
}

func getEnv(key, defaultValue string) string{
	if value := os.Getenv(key); value != ""{
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int{
	if value := os.Getenv(key); value != ""{
		if intValue, err := strconv.Atoi(value); err == nil{
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool{
	if value := os.Getenv(key); value != ""{
		if boolValue, err := strconv.ParseBool(value); err == nil{
			return boolValue
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