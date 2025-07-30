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

// TestRecommendationRepositoryCRUD tests basic CRUD operations for recommendations
func TestRecommendationRepositoryCRUD(t *testing.T) {
	// Setup test database
	testCfg := config.LoadTestConfig()
	if !testCfg.HasTestDatabase() {
		t.Skip("DATABASE_URL_TEST not set, skipping integration test")
	}

	// Connect to test database
	err := connectToTestDatabase()
	require.NoError(t, err)
	defer database.Close()

	// Create repository
	repo := NewRecommendationRepositorySimple(database.DB)

	// Test data
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
		// Clean up before test
		cleanupRecommendations(t, repo, now)

		// Create test stocks first (required for foreign key constraint)
		stockRepo := NewStockRepository(repo.GetDB())
		err := createTestStocksForRecommendations(t, stockRepo, testRecommendations)
		require.NoError(t, err)

		// Create recommendations
		err = repo.BulkCreate(testRecommendations)
		require.NoError(t, err)

		// Verify recommendations were created
		recommendations, err := repo.GetLatest(10)
		require.NoError(t, err)
		require.Len(t, recommendations, 3)

		// Verify data integrity
		for i, rec := range recommendations {
			assert.Equal(t, testRecommendations[i].Ticker, rec.Ticker)
			assert.Equal(t, testRecommendations[i].Score, rec.Score)
			assert.Equal(t, testRecommendations[i].Rank, rec.Rank)
			assert.Equal(t, testRecommendations[i].Explanation, rec.Explanation)
		}
	})

	t.Run("GetLatest Recommendations", func(t *testing.T) {
		// Clean up and recreate test data
		cleanupRecommendations(t, repo, now)
		err := repo.BulkCreate(testRecommendations)
		require.NoError(t, err)

		// Get latest recommendations with limit
		recommendations, err := repo.GetLatest(2)
		require.NoError(t, err)
		assert.Len(t, recommendations, 2)

		// Verify they are ordered by rank
		assert.Equal(t, 1, recommendations[0].Rank)
		assert.Equal(t, 2, recommendations[1].Rank)
	})

	t.Run("GetLatestRunAt", func(t *testing.T) {
		// Clean up and recreate test data
		cleanupRecommendations(t, repo, now)
		err := repo.BulkCreate(testRecommendations)
		require.NoError(t, err)

		// Get latest run_at
		latestRunAt, err := repo.GetLatestRunAt()
		require.NoError(t, err)
		require.NotNil(t, latestRunAt)

		// Verify it matches our test data (ignore timezone differences)
		assert.Equal(t, now.Truncate(time.Second).UTC(), latestRunAt.Truncate(time.Second).UTC())
	})

	t.Run("ShouldCalculateToday", func(t *testing.T) {
		// Clean up
		cleanupRecommendations(t, repo, now)

		// Test when no recommendations exist
		shouldCalculate, err := repo.ShouldCalculateToday()
		require.NoError(t, err)
		assert.True(t, shouldCalculate)

		// Create test stocks first (required for foreign key constraint)
		stockRepo := NewStockRepository(repo.GetDB())
		err = createTestStocksForRecommendations(t, stockRepo, testRecommendations)
		require.NoError(t, err)

		// Create recommendations for today
		err = repo.BulkCreate(testRecommendations)
		require.NoError(t, err)

		// Test when recommendations exist for today
		shouldCalculate, err = repo.ShouldCalculateToday()
		require.NoError(t, err)
		assert.False(t, shouldCalculate)

		// Clean up today's recommendations and create recommendations for yesterday
		cleanupRecommendations(t, repo, now)
		yesterday := now.Add(-24 * time.Hour)
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

		// Clean up and create old recommendations
		cleanupRecommendations(t, repo, yesterday)

		// Create test stocks for old recommendations
		stockRepo = NewStockRepository(repo.GetDB())
		err2 := createTestStocksForRecommendations(t, stockRepo, oldRecommendations)
		require.NoError(t, err2)

		err = repo.BulkCreate(oldRecommendations)
		require.NoError(t, err)

		// Test when only old recommendations exist
		shouldCalculate, err = repo.ShouldCalculateToday()
		require.NoError(t, err)
		assert.True(t, shouldCalculate)
	})

	t.Run("Empty Recommendations", func(t *testing.T) {
		// Clean up
		cleanupRecommendations(t, repo, now)

		// Test GetLatestRunAt with no recommendations first
		latestRunAt, err := repo.GetLatestRunAt()
		require.Error(t, err)
		assert.Nil(t, latestRunAt)

		// Test GetLatest with no recommendations
		recommendations, err := repo.GetLatest(10)
		require.Error(t, err) // Should return error when no recommendations exist
		assert.Empty(t, recommendations)
	})

	t.Run("BulkCreate Empty Slice", func(t *testing.T) {
		// Test with empty slice
		err := repo.BulkCreate([]*model.Recommendation{})
		require.NoError(t, err) // Should not error with empty slice
	})

	// Cleanup
	cleanupRecommendations(t, repo, now)
	cleanupRecommendations(t, repo, now.Add(-24*time.Hour))
}

// TestRecommendationRepositoryIntegration tests integration scenarios
func TestRecommendationRepositoryIntegration(t *testing.T) {
	// Setup test database
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

		// First run
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

		err := repo.BulkCreate(recommendations1)
		require.NoError(t, err)

		// Second run (should be the latest)
		recommendations2 := []*model.Recommendation{
			{
				ID:          uuid.New().String(),
				Ticker:      "TEST2",
				Score:       90.0,
				Explanation: "Second run",
				RunAt:       now.Add(time.Minute), // Slightly later
				Rank:        1,
			},
		}

		err = repo.BulkCreate(recommendations2)
		require.NoError(t, err)

		// GetLatest should return the second run
		latest, err := repo.GetLatest(10)
		require.NoError(t, err)
		require.Len(t, latest, 1)
		assert.Equal(t, "TEST2", latest[0].Ticker)
		assert.Equal(t, 90.0, latest[0].Score)
	})

	t.Run("Large Dataset", func(t *testing.T) {
		now := time.Now()
		cleanupRecommendations(t, repo, now)

		// Create many recommendations
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

		// Create test stocks first (required for foreign key constraint)
		stockRepo := NewStockRepository(repo.GetDB())
		err := createTestStocksForRecommendations(t, stockRepo, recommendations)
		require.NoError(t, err)

		err = repo.BulkCreate(recommendations)
		require.NoError(t, err)

		// Get latest with limit
		latest, err := repo.GetLatest(50)
		require.NoError(t, err)
		assert.Len(t, latest, 50)

		// Verify they are ordered by rank
		for i, rec := range latest {
			assert.Equal(t, i+1, rec.Rank)
		}
	})

	// Cleanup
	cleanupRecommendations(t, repo, time.Now())
}
