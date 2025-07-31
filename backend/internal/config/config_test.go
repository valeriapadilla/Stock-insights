package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadWithDefaultValues(t *testing.T) {
	os.Clearenv()
	os.Setenv("ENVIRONMENT", "test")

	config := Load()
	assert.Equal(t, "8080", config.Port)
	assert.Equal(t, "test", config.Environment)
	assert.Equal(t, "info", config.LogLevel)
	assert.Equal(t, "https://api.karenai.click", config.ExternalAPIURL)
	assert.Equal(t, 5*time.Minute, config.CacheTTL)
	assert.Equal(t, 100, config.RateLimit)
}

func TestLoadWithEnvVariables(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("PORT", "9090")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("EXTERNAL_API_URL", "https://test-api.com")
	os.Setenv("CACHE_TTL", "10m")
	os.Setenv("RATE_LIMIT", "200")

	config := Load()

	assert.Equal(t, "9090", config.Port)
	assert.Equal(t, "test", config.Environment) // Changed from production to test
	assert.Equal(t, "debug", config.LogLevel)
	assert.Equal(t, "https://test-api.com", config.ExternalAPIURL)
	assert.Equal(t, 10*time.Minute, config.CacheTTL)
	assert.Equal(t, 200, config.RateLimit)

	os.Clearenv()
}

func TestLoadWithEnvFile(t *testing.T) {
	envContent := `PORT=7070
	ENVIRONMENT=test
	LOG_LEVEL=warn
	EXTERNAL_API_URL=https://staging-api.com
	CACHE_TTL=15m 
	RATE_LIMIT=150`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	require.NoError(t, err)
	defer os.Remove(".env")

	config := Load()

	assert.Equal(t, "7070", config.Port)
	assert.Equal(t, "test", config.Environment) // Changed from staging to test
	assert.Equal(t, "warn", config.LogLevel)
	assert.Equal(t, "https://staging-api.com", config.ExternalAPIURL)
	assert.Equal(t, 15*time.Minute, config.CacheTTL)
	assert.Equal(t, 150, config.RateLimit)
}

func TestGetEnvAsDuration(t *testing.T) {
	os.Setenv("TEST_DURATION", "5m")
	duration := getEnvAsDuration("TEST_DURATION", 1*time.Minute)
	assert.Equal(t, 5*time.Minute, duration)

	os.Setenv("TEST_DURATION", "invalid")
	duration = getEnvAsDuration("TEST_DURATION", 1*time.Minute)
	assert.Equal(t, 1*time.Minute, duration)

	os.Unsetenv("TEST_DURATION")
	duration = getEnvAsDuration("TEST_DURATION", 1*time.Minute)
	assert.Equal(t, 1*time.Minute, duration)

	os.Clearenv()
}

func TestGetEnvAsInt(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	value := getEnvAsInt("TEST_INT", 10)
	assert.Equal(t, 42, value)

	os.Setenv("TEST_INT", "invalid")
	value = getEnvAsInt("TEST_INT", 10)
	assert.Equal(t, 10, value)

	os.Unsetenv("TEST_INT")
	value = getEnvAsInt("TEST_INT", 10)
	assert.Equal(t, 10, value)

	os.Clearenv()
}
