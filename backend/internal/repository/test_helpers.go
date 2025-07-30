package repository

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/valeriapadilla/stock-insights/internal/database"
	"github.com/valeriapadilla/stock-insights/internal/model"
)

// Helper function to connect to test database
func connectToTestDatabase() error {
	os.Setenv("TESTING", "true")
	defer os.Unsetenv("TESTING")

	return database.Connect()
}

// Helper function to create test stocks for recommendations
func createTestStocksForRecommendations(t *testing.T, repo *StockRepository, recommendations []*model.Recommendation) error {
	for _, rec := range recommendations {
		// Create a test stock for each recommendation
		testStock := &model.Stock{
			Ticker:     rec.Ticker,
			Company:    fmt.Sprintf("Test Company %s", rec.Ticker),
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

		// Clean up any existing stock first
		cleanupStock(t, repo, testStock.Ticker)

		// Create the stock
		err := createTestStock(repo, testStock)
		if err != nil {
			return fmt.Errorf("failed to create test stock for %s: %w", rec.Ticker, err)
		}
	}

	return nil
}

// Helper function to create test stock
func createTestStock(repo *StockRepository, stock *model.Stock) error {
	query := `
		INSERT INTO stocks (
			ticker, company, target_from, target_to, rating_from, rating_to,
			action, brokerage, time, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := repo.GetDB().Exec(query,
		stock.Ticker, stock.Company, stock.TargetFrom, stock.TargetTo,
		stock.RatingFrom, stock.RatingTo, stock.Action, stock.Brokerage,
		stock.Time, stock.CreatedAt, stock.UpdatedAt,
	)

	return err
}

// Helper function to cleanup stock
func cleanupStock(t *testing.T, repo *StockRepository, ticker string) {
	// First delete recommendations that reference this stock
	query := "DELETE FROM recommendations WHERE ticker = $1"
	_, err := repo.GetDB().Exec(query, ticker)
	require.NoError(t, err)

	// Then delete the stock
	query = "DELETE FROM stocks WHERE ticker = $1"
	_, err = repo.GetDB().Exec(query, ticker)
	require.NoError(t, err)
}

// Helper function to cleanup recommendations
func cleanupRecommendations(t *testing.T, repo *RecommendationRepositorySimple, runAt time.Time) {
	// Delete all recommendations to ensure clean state
	query := "DELETE FROM recommendations"
	_, err := repo.GetDB().Exec(query)
	require.NoError(t, err)
}

// Helper function to cleanup all test data
func cleanupAllTestData(t *testing.T, repo *StockRepository) {
	// Delete all recommendations first (due to foreign key constraint)
	query := "DELETE FROM recommendations"
	_, err := repo.GetDB().Exec(query)
	require.NoError(t, err)

	// Then delete all stocks
	query = "DELETE FROM stocks"
	_, err = repo.GetDB().Exec(query)
	require.NoError(t, err)

	// Verify cleanup was successful
	count, err := repo.Count(map[string]string{})
	require.NoError(t, err)
	require.Equal(t, 0, count, "Database should be empty after cleanup")
}
