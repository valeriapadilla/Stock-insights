package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valeriapadilla/stock-insights/internal/config"
	"github.com/valeriapadilla/stock-insights/internal/database"
	"github.com/valeriapadilla/stock-insights/internal/model"
)

// TestStockRepositoryCRUD tests basic CRUD operations for stocks
func TestStockRepositoryCRUD(t *testing.T) {
	err := connectToTestDatabase()
	require.NoError(t, err)
	defer database.Close()

	// Create repository
	repo := NewStockRepository(database.DB)

	// Test data
	testStock := &model.Stock{
		Ticker:     "TEST",
		Company:    "Test Company",
		TargetFrom: "$10.00",
		TargetTo:   "$15.00",
		RatingFrom: "Hold",
		RatingTo:   "Buy",
		Action:     "upgraded by",
		Brokerage:  "Test Brokerage",
		Time:       time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	t.Run("Create and Get Stock", func(t *testing.T) {
		// Clean up before test
		cleanupStock(t, repo, testStock.Ticker)

		// Create stock
		err := createTestStock(repo, testStock)
		require.NoError(t, err)

		// Get all stocks and verify
		stocks, err := repo.GetAll(10, 0, map[string]string{})
		require.NoError(t, err)
		require.NotEmpty(t, stocks)

		// Find our test stock
		var foundStock *model.Stock
		for _, stock := range stocks {
			if stock.Ticker == testStock.Ticker {
				foundStock = stock
				break
			}
		}

		require.NotNil(t, foundStock)
		assert.Equal(t, testStock.Ticker, foundStock.Ticker)
		assert.Equal(t, testStock.Company, foundStock.Company)
		assert.Equal(t, testStock.Brokerage, foundStock.Brokerage)
		assert.Equal(t, testStock.RatingTo, foundStock.RatingTo)
	})

	t.Run("Get Stocks with Filters", func(t *testing.T) {
		// Test filtering by brokerage
		filters := map[string]string{"brokerage": testStock.Brokerage}
		stocks, err := repo.GetAll(10, 0, filters)
		require.NoError(t, err)

		// Verify all returned stocks have the correct brokerage
		for _, stock := range stocks {
			assert.Equal(t, testStock.Brokerage, stock.Brokerage)
		}

		// Test filtering by rating
		filters = map[string]string{"rating": testStock.RatingTo}
		stocks, err = repo.GetAll(10, 0, filters)
		require.NoError(t, err)

		// Verify all returned stocks have the correct rating
		for _, stock := range stocks {
			assert.Equal(t, testStock.RatingTo, stock.RatingTo)
		}
	})

	t.Run("Get Stocks with Pagination", func(t *testing.T) {
		// Create multiple test stocks
		testStocks := []*model.Stock{
			{
				Ticker:     "TEST1",
				Company:    "Test Company 1",
				TargetFrom: "$10.00",
				TargetTo:   "$15.00",
				RatingFrom: "Hold",
				RatingTo:   "Buy",
				Action:     "upgraded by",
				Brokerage:  "Test Brokerage",
				Time:       time.Now(),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			{
				Ticker:     "TEST2",
				Company:    "Test Company 2",
				TargetFrom: "$20.00",
				TargetTo:   "$25.00",
				RatingFrom: "Sell",
				RatingTo:   "Hold",
				Action:     "downgraded by",
				Brokerage:  "Test Brokerage",
				Time:       time.Now(),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		}

		// Create test stocks
		for _, stock := range testStocks {
			cleanupStock(t, repo, stock.Ticker)
			err := createTestStock(repo, stock)
			require.NoError(t, err)
		}

		// Test pagination
		stocks, err := repo.GetAll(1, 0, map[string]string{})
		require.NoError(t, err)
		assert.Len(t, stocks, 1)

		// Test offset
		stocks, err = repo.GetAll(1, 1, map[string]string{})
		require.NoError(t, err)
		assert.Len(t, stocks, 1)
	})

	t.Run("Count Stocks", func(t *testing.T) {
		// Count all stocks
		count, err := repo.Count(map[string]string{})
		require.NoError(t, err)
		assert.Greater(t, count, 0)

		// Count with filter
		filters := map[string]string{"brokerage": testStock.Brokerage}
		count, err = repo.Count(filters)
		require.NoError(t, err)
		assert.Greater(t, count, 0)
	})

	// Cleanup
	cleanupStock(t, repo, testStock.Ticker)
	cleanupStock(t, repo, "TEST1")
	cleanupStock(t, repo, "TEST2")
}

// TestStockRepositoryIntegration tests integration scenarios
func TestStockRepositoryIntegration(t *testing.T) {
	// Setup test database
	testCfg := config.LoadTestConfig()
	if !testCfg.HasTestDatabase() {
		t.Skip("DATABASE_URL_TEST not set, skipping integration test")
	}

	err := connectToTestDatabase()
	require.NoError(t, err)
	defer database.Close()

	repo := NewStockRepository(database.DB)

	// Clean up any existing data before running integration tests
	cleanupAllTestData(t, repo)

	t.Run("Empty Database", func(t *testing.T) {
		// Clean up any existing data first
		cleanupAllTestData(t, repo)

		// Test with empty database
		stocks, err := repo.GetAll(10, 0, map[string]string{})
		require.NoError(t, err)
		assert.Empty(t, stocks)

		count, err := repo.Count(map[string]string{})
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("Invalid Filters", func(t *testing.T) {
		// Test with invalid filters (should not cause errors)
		filters := map[string]string{"invalid_field": "value"}
		_, err := repo.GetAll(10, 0, filters)
		require.NoError(t, err)
		// Should return empty result, not error
	})

	t.Run("Large Limit", func(t *testing.T) {
		// Test with large limit
		stocks, err := repo.GetAll(1000, 0, map[string]string{})
		require.NoError(t, err)
		// Should not cause memory issues
		assert.NotNil(t, stocks)
	})
}
