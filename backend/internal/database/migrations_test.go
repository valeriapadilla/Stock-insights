package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/valeriapadilla/stock-insights/internal/config"
)

func init() {
	os.Setenv("LOG_LEVEL", "error")
}

func TestRunMigrations(t *testing.T) {
	testConfig := config.LoadTestConfig()
	if !testConfig.HasTestDatabase() {
		t.Skip("DATABASE_URL_TEST not set for integration tests")
	}

	testDBURL := testConfig.GetTestDatabaseURL()

	os.Setenv("DATABASE_URL", testDBURL)
	defer os.Unsetenv("DATABASE_URL")

	err := Connect()
	if err != nil {
		t.Skipf("Database not available for testing: %v", err)
	}
	defer Close()

	cleanupTestDatabase(t)

	manager := NewMigrationManager(DB)

	err = manager.RunMigrations()
	require.NoError(t, err)

	verifyTablesExist(t)
	verifyIndexesExist(t)
}

func TestRunMigrationsWithoutConnection(t *testing.T) {
	manager := NewMigrationManager(nil)
	err := manager.RunMigrations()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database connection not established")
}

func TestMigrationsAreIdempotent(t *testing.T) {
	testConfig := config.LoadTestConfig()
	if !testConfig.HasTestDatabase() {
		t.Skip("DATABASE_URL_TEST not set for integration tests")
	}

	testDBURL := testConfig.GetTestDatabaseURL()

	os.Setenv("DATABASE_URL", testDBURL)
	defer os.Unsetenv("DATABASE_URL")

	err := Connect()
	if err != nil {
		t.Skipf("Database not available for testing: %v", err)
	}
	defer Close()

	cleanupTestDatabase(t)

	manager := NewMigrationManager(DB)

	err = manager.RunMigrations()
	require.NoError(t, err)

	err = manager.RunMigrations()
	require.NoError(t, err)

	verifyTablesExist(t)
}

func TestMigrationsTableStructure(t *testing.T) {
	testConfig := config.LoadTestConfig()
	if !testConfig.HasTestDatabase() {
		t.Skip("DATABASE_URL_TEST not set for integration tests")
	}

	testDBURL := testConfig.GetTestDatabaseURL()

	os.Setenv("DATABASE_URL", testDBURL)
	defer os.Unsetenv("DATABASE_URL")

	err := Connect()
	if err != nil {
		t.Skipf("Database not available for testing: %v", err)
	}
	defer Close()

	cleanupTestDatabase(t)

	manager := NewMigrationManager(DB)
	err = manager.RunMigrations()
	require.NoError(t, err)

	verifyMigrationsTableStructure(t)
}

func cleanupTestDatabase(t *testing.T) {
	queries := []string{
		"DROP TABLE IF EXISTS recommendations CASCADE",
		"DROP TABLE IF EXISTS stocks CASCADE",
		"DROP TABLE IF EXISTS migrations CASCADE",
	}

	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			t.Logf("Warning: Could not cleanup table: %v", err)
		}
	}
}

func verifyTablesExist(t *testing.T) {
	tables := []string{"stocks", "recommendations", "migrations"}

	for _, tableName := range tables {
		var exists bool
		err := DB.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM information_schema.tables 
				WHERE table_name = $1
			)
		`, tableName).Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists, "Table %s should exist", tableName)
	}
}

func verifyIndexesExist(t *testing.T) {
	indexes := []string{
		"idx_stocks_brokerage",
		"idx_stocks_rating_to",
		"idx_stocks_time",
		"idx_stocks_company",
		"idx_stocks_brokerage_time",
		"idx_recommendations_run_at",
		"idx_recommendations_score",
		"idx_recommendations_rank",
		"idx_recommendations_ticker",
		"idx_recommendations_run_at_score",
	}

	for _, indexName := range indexes {
		var exists bool
		err := DB.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM pg_indexes WHERE indexname = $1
			)
		`, indexName).Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists, "Index %s should exist", indexName)
	}
}

func verifyMigrationsTableStructure(t *testing.T) {
	var columnName, dataType string
	rows, err := DB.Query(`
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = 'migrations' 
		ORDER BY ordinal_position
	`)
	require.NoError(t, err)
	defer rows.Close()

	expectedColumns := map[string]string{
		"id":          "text",
		"filename":    "text",
		"executed_at": "timestamp with time zone",
	}

	foundColumns := make(map[string]string)
	for rows.Next() {
		err := rows.Scan(&columnName, &dataType)
		require.NoError(t, err)
		foundColumns[columnName] = dataType
	}

	for expectedCol, expectedType := range expectedColumns {
		assert.Contains(t, foundColumns, expectedCol)
		if foundType, exists := foundColumns[expectedCol]; exists {
			assert.Equal(t, expectedType, foundType)
		}
	}
}
