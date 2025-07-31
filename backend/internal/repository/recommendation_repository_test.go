package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valeriapadilla/stock-insights/internal/model"
)

func TestRecommendationRepositoryCRUD(t *testing.T) {
	repo, cleanup := setupRecommendationTest(t)
	defer cleanup()

	now := time.Now()
	testRecommendations := createTestRecommendationsWithCustomData(
		[]string{"TEST1", "TEST2", "TEST3"},
		[]float64{85.5, 82.3, 78.9},
		now,
	)

	t.Run("BulkCreate Recommendations", func(t *testing.T) {
		cleanupRecommendationTest(t, repo, now)

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
		cleanupRecommendationTest(t, repo, now)
		err := repo.BulkCreate(testRecommendations)
		require.NoError(t, err)

		recommendations, err := repo.GetLatest(2)
		require.NoError(t, err)
		assert.Len(t, recommendations, 2)

		assert.Equal(t, 1, recommendations[0].Rank)
		assert.Equal(t, 2, recommendations[1].Rank)
	})

	t.Run("GetLatestRunAt", func(t *testing.T) {
		cleanupRecommendationTest(t, repo, now)
		err := repo.BulkCreate(testRecommendations)
		require.NoError(t, err)

		latestRunAt, err := repo.GetLatestRunAt()
		require.NoError(t, err)
		require.NotNil(t, latestRunAt)

		assert.Equal(t, now.Truncate(time.Second).UTC(), latestRunAt.Truncate(time.Second).UTC())
	})

	t.Run("Empty Recommendations", func(t *testing.T) {
		_, err := repo.GetDB().Exec("DELETE FROM recommendations")
		require.NoError(t, err)

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

	cleanupRecommendationTest(t, repo, now)
	cleanupRecommendationTest(t, repo, now.Add(-24*time.Hour))
}

func TestRecommendationRepositoryIntegration(t *testing.T) {
	repo, cleanup := setupRecommendationTest(t)
	defer cleanup()

	t.Run("Multiple Runs Same Day", func(t *testing.T) {
		now := time.Now()
		cleanupRecommendationTest(t, repo, now)

		recommendations1 := createTestRecommendationsWithCustomData(
			[]string{"TEST1"},
			[]float64{85.5},
			now,
		)

		stockRepo := NewStockRepository(repo.GetDB())
		err := createTestStocksForRecommendations(t, stockRepo, recommendations1)
		require.NoError(t, err)

		err = repo.BulkCreate(recommendations1)
		require.NoError(t, err)

		recommendations2 := createTestRecommendationsWithCustomData(
			[]string{"TEST2"},
			[]float64{90.0},
			now.Add(time.Minute),
		)

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
		cleanupRecommendationTest(t, repo, now)

		recommendations := createTestRecommendations(100, now)

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

	cleanupRecommendationTest(t, repo, time.Now())
}
