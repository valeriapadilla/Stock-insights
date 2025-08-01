package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/valeriapadilla/stock-insights/internal/model"
	"github.com/valeriapadilla/stock-insights/internal/repository/interfaces"
	"github.com/valeriapadilla/stock-insights/internal/validator"
)

// RecommendationRepositorySimple implements interfaces.RecommendationRepository
type RecommendationRepositorySimple struct {
	*BaseRepository
	validator *validator.RecommendationValidator
}

// Ensure RecommendationRepositorySimple implements the interface
var _ interfaces.RecommendationRepository = (*RecommendationRepositorySimple)(nil)

func NewRecommendationRepositorySimple(db *sql.DB) *RecommendationRepositorySimple {
	return &RecommendationRepositorySimple{
		BaseRepository: NewBaseRepository(db),
		validator:      validator.NewRecommendationValidator(),
	}
}

func (r *RecommendationRepositorySimple) GetLatest(limit int) ([]*model.Recommendation, error) {
	if limit <= 0 || limit > 1000 {
		limit = 10
	}

	latestRunAt, err := r.GetLatestRunAt()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest run_at: %w", err)
	}

	qb := NewQueryBuilder().
		Select("id", "ticker", "score", "explanation", "run_at", "rank").
		From("recommendations").
		Where("run_at = $1", latestRunAt).
		OrderBy("rank", "ASC").
		Limit(limit)

	query, args := qb.Build()

	rows, err := r.GetDB().Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest recommendations: %w", err)
	}
	defer rows.Close()

	var recommendations []*model.Recommendation
	for rows.Next() {
		var recommendation model.Recommendation
		err := rows.Scan(
			&recommendation.ID, &recommendation.Ticker, &recommendation.Score,
			&recommendation.Explanation, &recommendation.RunAt, &recommendation.Rank,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recommendation: %w", err)
		}
		recommendations = append(recommendations, &recommendation)
	}

	if len(recommendations) == 0 {
		return nil, fmt.Errorf("no recommendations found")
	}

	return recommendations, nil
}

func (r *RecommendationRepositorySimple) GetLatestRunAt() (*time.Time, error) {
	query := "SELECT MAX(run_at) FROM recommendations"

	var runAt *time.Time
	err := r.GetDB().QueryRow(query).Scan(&runAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no recommendations found")
		}
		return nil, fmt.Errorf("failed to get latest run_at: %w", err)
	}

	if runAt == nil {
		return nil, fmt.Errorf("no recommendations found")
	}

	return runAt, nil
}
