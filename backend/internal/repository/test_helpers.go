package repository

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/valeriapadilla/stock-insights/internal/config"
	"github.com/valeriapadilla/stock-insights/internal/database"
	"github.com/valeriapadilla/stock-insights/internal/model"
)

func connectToTestDatabase() error {
	os.Setenv("TESTING", "true")
	defer os.Unsetenv("TESTING")

	err := database.Connect()
	if err != nil {
		return err
	}

	cleanDatabase()

	manager := database.NewMigrationManager(database.DB)
	return manager.RunMigrations()
}

func cleanDatabase() {
	queries := []string{
		"DROP TABLE IF EXISTS recommendations CASCADE",
		"DROP TABLE IF EXISTS stocks CASCADE",
		"DELETE FROM migrations",
		"DROP TABLE IF EXISTS migrations CASCADE",
	}

	for _, query := range queries {
		_, err := database.DB.Exec(query)
		if err != nil {
			fmt.Printf("Warning: failed to clean database: %v\n", err)
		}
	}
}

func setupRecommendationTest(t *testing.T) (*RecommendationRepositorySimple, func()) {
	testCfg := config.LoadTestConfig()
	if !testCfg.HasTestDatabase() {
		t.Skip("DATABASE_URL_TEST not set, skipping integration test")
	}

	err := connectToTestDatabase()
	require.NoError(t, err)

	repo := NewRecommendationRepositorySimple(database.DB)

	cleanup := func() {
		database.Close()
	}

	return repo, cleanup
}

func createTestRecommendations(count int, baseTime time.Time) []*model.Recommendation {
	recommendations := make([]*model.Recommendation, count)
	for i := 0; i < count; i++ {
		recommendations[i] = &model.Recommendation{
			ID:          uuid.New().String(),
			Ticker:      fmt.Sprintf("TEST%d", i),
			Score:       float64(100 - i),
			Explanation: fmt.Sprintf("Test recommendation %d", i),
			RunAt:       baseTime,
			Rank:        i + 1,
		}
	}
	return recommendations
}

func createTestRecommendationsWithCustomData(tickers []string, scores []float64, baseTime time.Time) []*model.Recommendation {
	recommendations := make([]*model.Recommendation, len(tickers))
	for i := 0; i < len(tickers); i++ {
		score := 85.0
		if i < len(scores) {
			score = scores[i]
		}
		recommendations[i] = &model.Recommendation{
			ID:          uuid.New().String(),
			Ticker:      tickers[i],
			Score:       score,
			Explanation: fmt.Sprintf("Test recommendation for %s", tickers[i]),
			RunAt:       baseTime,
			Rank:        i + 1,
		}
	}
	return recommendations
}

func cleanupRecommendationTest(t *testing.T, repo *RecommendationRepositorySimple, runAt time.Time) {
	cleanupRecommendations(t, repo, runAt)
}

func createTestStocksForRecommendations(t *testing.T, repo *StockRepository, recommendations []*model.Recommendation) error {
	for _, rec := range recommendations {
		if rec.Ticker == "" || len(rec.Ticker) > 10 {
			return fmt.Errorf("invalid ticker in test data: %s", rec.Ticker)
		}

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

		cleanupStock(t, repo, testStock.Ticker)

		err := createTestStock(repo, testStock)
		if err != nil {
			return fmt.Errorf("failed to create test stock for %s: %w", rec.Ticker, err)
		}
	}

	return nil
}

func createTestStock(repo *StockRepository, stock *model.Stock) error {
	if stock.Ticker == "" || len(stock.Ticker) > 10 {
		return fmt.Errorf("invalid ticker")
	}

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

func cleanupStock(t *testing.T, repo *StockRepository, ticker string) {
	if ticker == "" || len(ticker) > 10 {
		t.Fatalf("invalid ticker for cleanup: %s", ticker)
	}

	query := "DELETE FROM recommendations WHERE ticker = $1"
	_, err := repo.GetDB().Exec(query, ticker)
	require.NoError(t, err)

	query = "DELETE FROM stocks WHERE ticker = $1"
	_, err = repo.GetDB().Exec(query, ticker)
	require.NoError(t, err)
}

func cleanupRecommendations(t *testing.T, repo *RecommendationRepositorySimple, runAt time.Time) {
	query := "DELETE FROM recommendations WHERE run_at::date = $1::date"
	_, err := repo.GetDB().Exec(query, runAt)
	require.NoError(t, err)
}

func cleanupAllTestData(t *testing.T, repo *StockRepository) {
	query := "DELETE FROM recommendations"
	_, err := repo.GetDB().Exec(query)
	require.NoError(t, err)

	query = "DELETE FROM stocks"
	_, err = repo.GetDB().Exec(query)
	require.NoError(t, err)

	count, err := repo.Count(map[string]string{})
	require.NoError(t, err)
	require.Equal(t, 0, count, "Database should be empty after cleanup")
}
