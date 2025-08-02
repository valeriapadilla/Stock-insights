package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valeriapadilla/stock-insights/internal/config"
	"github.com/valeriapadilla/stock-insights/internal/database"
	"github.com/valeriapadilla/stock-insights/internal/model"
	repoInterfaces "github.com/valeriapadilla/stock-insights/internal/repository/interfaces"
)

func TestStockRepositoryCRUD(t *testing.T) {
	testCfg := config.LoadTestConfig()
	if !testCfg.HasTestDatabase() {
		t.Skip("DATABASE_URL_TEST not set, skipping integration test")
	}

	err := connectToTestDatabase()
	require.NoError(t, err)
	defer database.Close()

	repo := NewStockRepository(database.DB)

	cleanupAllTestData(t, repo)

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
		cleanupStock(t, repo, testStock.Ticker)

		err := createTestStock(repo, testStock)
		require.NoError(t, err)

		stocks, err := repo.GetStocks(repoInterfaces.GetStocksParams{
			Limit:  10,
			Offset: 0,
			Sort:   "time",
			Order:  "desc",
		})
		require.NoError(t, err)
		require.NotEmpty(t, stocks)

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

	t.Run("Get Stocks", func(t *testing.T) {
		cleanupStock(t, repo, "TEST1")
		cleanupStock(t, repo, "TEST2")
		cleanupStock(t, repo, "TEST3")

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
			{
				Ticker:     "TEST3",
				Company:    "Test Company 3",
				TargetFrom: "$30.00",
				TargetTo:   "$35.00",
				RatingFrom: "Hold",
				RatingTo:   "Buy",
				Action:     "BUY",
				Brokerage:  "Test Brokerage",
				Time:       time.Now(),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		}

		for _, stock := range testStocks {
			err := createTestStock(repo, stock)
			require.NoError(t, err)
		}

		stocks, err := repo.GetStocks(repoInterfaces.GetStocksParams{
			Limit:  10,
			Offset: 0,
			Sort:   "time",
			Order:  "desc",
		})
		require.NoError(t, err)
		assert.NotNil(t, stocks)
		assert.GreaterOrEqual(t, len(stocks), 3)
	})

	t.Run("Get Stocks With Filters", func(t *testing.T) {
		filters := map[string]string{
			"brokerage": "Test Brokerage",
		}

		stocks, err := repo.GetStocks(repoInterfaces.GetStocksParams{
			Limit:   10,
			Offset:  0,
			Sort:    "time",
			Order:   "desc",
			Filters: filters,
		})
		require.NoError(t, err)
		assert.NotNil(t, stocks)
		assert.GreaterOrEqual(t, len(stocks), 1)
		if len(stocks) > 0 {
			assert.Equal(t, "Test Brokerage", stocks[0].Brokerage)
		}
	})

	t.Run("Get Stocks With Multiple Filters", func(t *testing.T) {
		filters := map[string]string{
			"brokerage": "Test Brokerage",
			"action":    "BUY",
		}

		stocks, err := repo.GetStocks(repoInterfaces.GetStocksParams{
			Limit:   10,
			Offset:  0,
			Sort:    "time",
			Order:   "desc",
			Filters: filters,
		})
		require.NoError(t, err)
		assert.NotNil(t, stocks)
		assert.Len(t, stocks, 1)
		assert.Equal(t, "Test Brokerage", stocks[0].Brokerage)
		assert.Equal(t, "BUY", stocks[0].Action)
	})

	t.Run("Get Stocks with Pagination", func(t *testing.T) {
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

		for _, stock := range testStocks {
			cleanupStock(t, repo, stock.Ticker)
			err := createTestStock(repo, stock)
			require.NoError(t, err)
		}

		stocks, err := repo.GetStocks(repoInterfaces.GetStocksParams{
			Limit:  1,
			Offset: 0,
			Sort:   "time",
			Order:  "desc",
		})
		require.NoError(t, err)
		assert.Len(t, stocks, 1)

		stocks, err = repo.GetStocks(repoInterfaces.GetStocksParams{
			Limit:  1,
			Offset: 1,
			Sort:   "time",
			Order:  "desc",
		})
		require.NoError(t, err)
		assert.Len(t, stocks, 1)
	})

	t.Run("Count Stocks", func(t *testing.T) {
		count, err := repo.GetStocksCount(repoInterfaces.GetStocksParams{})
		require.NoError(t, err)
		assert.Greater(t, count, 0)

		filters := map[string]string{"brokerage": testStock.Brokerage}
		count, err = repo.GetStocksCount(repoInterfaces.GetStocksParams{
			Filters: filters,
		})
		require.NoError(t, err)
		assert.Greater(t, count, 0)
	})

	cleanupStock(t, repo, testStock.Ticker)
	cleanupStock(t, repo, "TEST1")
	cleanupStock(t, repo, "TEST2")
}

func TestStockRepositoryIntegration(t *testing.T) {
	testCfg := config.LoadTestConfig()
	if !testCfg.HasTestDatabase() {
		t.Skip("DATABASE_URL_TEST not set, skipping integration test")
	}

	err := connectToTestDatabase()
	require.NoError(t, err)
	defer database.Close()

	repo := NewStockRepository(database.DB)

	cleanupAllTestData(t, repo)

	t.Run("Empty Database", func(t *testing.T) {
		cleanupAllTestData(t, repo)

		stocks, err := repo.GetStocks(repoInterfaces.GetStocksParams{
			Limit:  10,
			Offset: 0,
			Sort:   "time",
			Order:  "desc",
		})
		require.NoError(t, err)
		assert.Empty(t, stocks)

		count, err := repo.GetStocksCount(repoInterfaces.GetStocksParams{})
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("Invalid Filters", func(t *testing.T) {
		filters := map[string]string{"invalid_field": "value"}
		_, err := repo.GetStocks(repoInterfaces.GetStocksParams{
			Limit:   10,
			Offset:  0,
			Sort:    "time",
			Order:   "desc",
			Filters: filters,
		})
		require.NoError(t, err)
	})

	t.Run("Large Limit", func(t *testing.T) {
		stocks, err := repo.GetStocks(repoInterfaces.GetStocksParams{
			Limit:  1000,
			Offset: 0,
			Sort:   "time",
			Order:  "desc",
		})
		require.NoError(t, err)
		_ = stocks
	})
}
