package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	config := Load()

	assert.NotNil(t, config)
	assert.Equal(t, "8080", config.Port)
	assert.Equal(t, "development", config.Environment)
	assert.Equal(t, "info", config.LogLevel)
	assert.Equal(t, "https://api.karenai.click", config.ExternalAPIURL)
	assert.Equal(t, 5*time.Minute, config.CacheTTL)
	assert.Equal(t, 100, config.RateLimit)
}

func TestConfig_LoadWithEnvironment(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("DATABASE_URL", "postgres://test")
	os.Setenv("EXTERNAL_API_KEY", "test-key")
	os.Setenv("CACHE_TTL", "10m")
	os.Setenv("RATE_LIMIT", "200")
	os.Setenv("ADMIN_API_KEY", "admin-key")

	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("ENVIRONMENT")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("EXTERNAL_API_KEY")
		os.Unsetenv("CACHE_TTL")
		os.Unsetenv("RATE_LIMIT")
		os.Unsetenv("ADMIN_API_KEY")
	}()

	config := Load()

	assert.Equal(t, "9090", config.Port)
	assert.Equal(t, "production", config.Environment)
	assert.Equal(t, "debug", config.LogLevel)
	assert.Equal(t, "postgres://test", config.DatabaseURL)
	assert.Equal(t, "test-key", config.ExternalAPIKey)
	assert.Equal(t, 10*time.Minute, config.CacheTTL)
	assert.Equal(t, 200, config.RateLimit)
	assert.Equal(t, "admin-key", config.AdminAPIKey)
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "valid config",
			config: &Config{
				Environment: "development",
				DatabaseURL: "postgres://test",
				RateLimit:   100,
			},
			expectError: false,
		},
		{
			name: "test environment - no validation",
			config: &Config{
				Environment: "test",
				DatabaseURL: "",
				RateLimit:   0,
			},
			expectError: false,
		},
		{
			name: "missing database URL",
			config: &Config{
				Environment: "development",
				DatabaseURL: "",
				RateLimit:   100,
			},
			expectError: true,
		},
		{
			name: "invalid environment",
			config: &Config{
				Environment: "invalid",
				DatabaseURL: "postgres://test",
				RateLimit:   100,
			},
			expectError: true,
		},
		{
			name: "invalid rate limit",
			config: &Config{
				Environment: "development",
				DatabaseURL: "postgres://test",
				RateLimit:   0,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "test-value")
	defer os.Unsetenv("TEST_KEY")

	value := getEnv("TEST_KEY", "default")
	assert.Equal(t, "test-value", value)

	value = getEnv("NON_EXISTENT_KEY", "default")
	assert.Equal(t, "default", value)
}

func TestGetEnvAsInt(t *testing.T) {
	os.Setenv("TEST_INT", "123")
	defer os.Unsetenv("TEST_INT")

	value := getEnvAsInt("TEST_INT", 0)
	assert.Equal(t, 123, value)

	os.Setenv("TEST_INVALID", "not-a-number")
	defer os.Unsetenv("TEST_INVALID")

	value = getEnvAsInt("TEST_INVALID", 456)
	assert.Equal(t, 456, value)

	value = getEnvAsInt("NON_EXISTENT_INT", 789)
	assert.Equal(t, 789, value)
}

func TestGetEnvAsDuration(t *testing.T) {
	os.Setenv("TEST_DURATION", "5m")
	defer os.Unsetenv("TEST_DURATION")

	value := getEnvAsDuration("TEST_DURATION", time.Hour)
	assert.Equal(t, 5*time.Minute, value)

	os.Setenv("TEST_INVALID_DURATION", "not-a-duration")
	defer os.Unsetenv("TEST_INVALID_DURATION")

	value = getEnvAsDuration("TEST_INVALID_DURATION", time.Hour)
	assert.Equal(t, time.Hour, value)

	value = getEnvAsDuration("NON_EXISTENT_DURATION", 2*time.Hour)
	assert.Equal(t, 2*time.Hour, value)
}
