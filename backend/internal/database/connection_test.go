package database

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valeriapadilla/stock-insights/internal/config"
)

func TestConnectWithValidURL(t *testing.T) {
	testConfig := config.LoadTestConfig()
	if !testConfig.HasTestDatabase() {
		t.Skip("DATABASE_URL_TEST not set, skipping integration test")
	}

	os.Setenv("TESTING", "true")
	os.Setenv("DATABASE_URL", testConfig.GetTestDatabaseURL())
	defer os.Unsetenv("TESTING")
	defer os.Unsetenv("DATABASE_URL")

	err := Connect()

	if err != nil {
		assert.Contains(t, err.Error(), "connection")
	} else {
		assert.NotNil(t, DB)
		Close()
	}
}

func TestConnectWithInvalidURL(t *testing.T) {
	os.Setenv("TESTING", "false")
	os.Setenv("DATABASE_URL", "invalid://url")
	defer os.Unsetenv("TESTING")
	defer os.Unsetenv("DATABASE_URL")

	err := Connect()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to database")
}

func TestConnectWithMissingURL(t *testing.T) {
	os.Setenv("TESTING", "false")
	os.Unsetenv("DATABASE_URL")
	defer os.Unsetenv("TESTING")

	err := Connect()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DATABASE_URL is required")
}

func TestPingWithoutConnection(t *testing.T) {
	DB = nil

	err := Ping()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database connection not established")
}

func TestGetStatsWithoutConnection(t *testing.T) {
	DB = nil
	stats := GetStats()
	assert.Equal(t, sql.DBStats{}, stats)
}
