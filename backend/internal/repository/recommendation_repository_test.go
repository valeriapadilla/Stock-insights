package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valeriapadilla/stock-insights/internal/config"
	"github.com/valeriapadilla/stock-insights/internal/database"
	"github.com/valeriapadilla/stock-insights/internal/model"
)

func TestRecommendationRepositoryCRUD(t *testing.T) {
	testCfg := config.LoadTestConfig()
	if !testCfg.HasTestDatabase() {
		t.Skip("DATABASE_URL_TEST not set, skipping integration test")
	}

	err := connectToTestDatabase()
	require.NoError(t, err)
	defer database.Close()

	repo := NewRecommendationRepositorySimple(database.DB)

	now := time.Now()
	testRecommendations := []*model.Recommendation{
		{
			ID:          uuid.New().String(),
			Ticker:      "TEST1",
			Score:       85.5,
			Explanation: "Strong buy rating from HC Wainwright",
			RunAt:       now,
			Rank:        1,
		},
		{
			ID:          uuid.New().String(),
			Ticker:      "TEST2",
			Score:       82.3,
			Explanation: "Upgraded to Buy by Test Brokerage",
			RunAt:       now,
			Rank:        2,
		},
		{
			ID:          uuid.New().String(),
			Ticker:      "TEST3",
			Score:       78.9,
			Explanation: "Maintained Hold rating",
			RunAt:       now,
			Rank:        3,
		},
	}

	t.Run("BulkCreate Recommendations", func(t *testing.T) {
		cleanupRecommendations(t, repo, now)

		stockRepo := NewStockRepository(repo.GetDB())
		err := createTestStocksForRecommendations(t, stockRepo, testRecommendations)
		require.NoError(t, err)

		err = repo.BulkCreate(testRecommendations)
		require.NoError(t, err)

		recommendations, err := repo.GetLatest(10)
		require.NoError(t, err)
		require.Len(t, recommendations, 3)

		for i, rec := range recommendations {
			assert.Equal(t, testRecommendations[i].Ticker, rec.Ticker)
			assert.Equal(t, testRecommendations[i].Score, rec.Score)
			assert.Equal(t, testRecommendations[i].Rank, rec.Rank)
			assert.Equal(t, testRecommendations[i].Explanation, rec.Explanation)
		}
	})

	t.Run("GetLatest Recommendations", func(t *testing.T) {
		cleanupRecommendations(t, repo, now)
		err := repo.BulkCreate(testRecommendations)
		require.NoError(t, err)

		recommendations, err := repo.GetLatest(2)
		require.NoError(t, err)
		assert.Len(t, recommendations, 2)

		assert.Equal(t, 1, recommendations[0].Rank)
		assert.Equal(t, 2, recommendations[1].Rank)
	})

	t.Run("GetLatestRunAt", func(t *testing.T) {
		cleanupRecommendations(t, repo, now)
		err := repo.BulkCreate(testRecommendations)
		require.NoError(t, err)

		latestRunAt, err := repo.GetLatestRunAt()
		require.NoError(t, err)
		require.NotNil(t, latestRunAt)

		assert.Equal(t, now.Truncate(time.Second).UTC(), latestRunAt.Truncate(time.Second).UTC())
	})

	t.Run("ShouldCalculateToday", func(t *testing.T) {
		cleanupRecommendations(t, repo, now)

		shouldCalculate, err := repo.ShouldCalculateToday()
		require.NoError(t, err)
		assert.True(t, shouldCalculate)

		stockRepo := NewStockRepository(repo.GetDB())
		err = createTestStocksForRecommendations(t, stockRepo, testRecommendations)
		require.NoError(t, err)

		err = repo.BulkCreate(testRecommendations)
		require.NoError(t, err)

		shouldCalculate, err = repo.ShouldCalculateToday()
		require.NoError(t, err)
		assert.False(t, shouldCalculate)

		// Clean up today's recommendations before testing yesterday
		yesterday := now.Add(-24 * time.Hour)
		cleanupRecommendations(t, repo, now)

		oldRecommendations := []*model.Recommendation{
			{
				ID:          uuid.New().String(),
				Ticker:      "OLD1",
				Score:       75.0,
				Explanation: "Old recommendation",
				RunAt:       yesterday,
				Rank:        1,
			},
		}

		stockRepo = NewStockRepository(repo.GetDB())
		err2 := createTestStocksForRecommendations(t, stockRepo, oldRecommendations)
		require.NoError(t, err2)

		err = repo.BulkCreate(oldRecommendations)
		require.NoError(t, err)

		shouldCalculate, err = repo.ShouldCalculateToday()
		require.NoError(t, err)
		assert.True(t, shouldCalculate)
	})

	t.Run("Empty Recommendations", func(t *testing.T) {
		// Clean ALL recommendations, not just today's
		_, err := repo.GetDB().Exec("DELETE FROM recommendations")
		require.NoError(t, err)

		// Verify database is empty
		var count int
		err = repo.GetDB().QueryRow("SELECT COUNT(*) FROM recommendations").Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 0, count, "Database should be empty")

		latestRunAt, err := repo.GetLatestRunAt()
		require.Error(t, err)
		assert.Nil(t, latestRunAt)

		recommendations, err := repo.GetLatest(10)
		require.Error(t, err)
		assert.Empty(t, recommendations)
	})

	t.Run("BulkCreate Empty Slice", func(t *testing.T) {
		err := repo.BulkCreate([]*model.Recommendation{})
		require.NoError(t, err)
	})

	cleanupRecommendations(t, repo, now)
	cleanupRecommendations(t, repo, now.Add(-24*time.Hour))
}

func TestRecommendationRepositoryIntegration(t *testing.T) {
	testCfg := config.LoadTestConfig()
	if !testCfg.HasTestDatabase() {
		t.Skip("DATABASE_URL_TEST not set, skipping integration test")
	}

	err := connectToTestDatabase()
	require.NoError(t, err)
	defer database.Close()

	repo := NewRecommendationRepositorySimple(database.DB)

	t.Run("Multiple Runs Same Day", func(t *testing.T) {
		now := time.Now()
		cleanupRecommendations(t, repo, now)

		recommendations1 := []*model.Recommendation{
			{
				ID:          uuid.New().String(),
				Ticker:      "TEST1",
				Score:       85.5,
				Explanation: "First run",
				RunAt:       now,
				Rank:        1,
			},
		}

		stockRepo := NewStockRepository(repo.GetDB())
		err := createTestStocksForRecommendations(t, stockRepo, recommendations1)
		require.NoError(t, err)

		err = repo.BulkCreate(recommendations1)
		require.NoError(t, err)

		recommendations2 := []*model.Recommendation{
			{
				ID:          uuid.New().String(),
				Ticker:      "TEST2",
				Score:       90.0,
				Explanation: "Second run",
				RunAt:       now.Add(time.Minute),
				Rank:        1,
			},
		}

		err = createTestStocksForRecommendations(t, stockRepo, recommendations2)
		require.NoError(t, err)

		err = repo.BulkCreate(recommendations2)
		require.NoError(t, err)

		latest, err := repo.GetLatest(10)
		require.NoError(t, err)
		require.Len(t, latest, 1)
		assert.Equal(t, "TEST2", latest[0].Ticker)
		assert.Equal(t, 90.0, latest[0].Score)
	})

	t.Run("Large Dataset", func(t *testing.T) {
		now := time.Now()
		cleanupRecommendations(t, repo, now)

		var recommendations []*model.Recommendation
		for i := 0; i < 100; i++ {
			recommendations = append(recommendations, &model.Recommendation{
				ID:          uuid.New().String(),
				Ticker:      fmt.Sprintf("TEST%d", i),
				Score:       float64(100 - i),
				Explanation: fmt.Sprintf("Test recommendation %d", i),
				RunAt:       now,
				Rank:        i + 1,
			})
		}

		stockRepo := NewStockRepository(repo.GetDB())
		err := createTestStocksForRecommendations(t, stockRepo, recommendations)
		require.NoError(t, err)

		err = repo.BulkCreate(recommendations)
		require.NoError(t, err)

		latest, err := repo.GetLatest(50)
		require.NoError(t, err)
		assert.Len(t, latest, 50)

		for i, rec := range latest {
			assert.Equal(t, i+1, rec.Rank)
		}
	})

	cleanupRecommendations(t, repo, time.Now())
}
