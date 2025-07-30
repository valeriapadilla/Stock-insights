package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/valeriapadilla/stock-insights/internal/model"
)

type RecommendationRepositorySimple struct {
	*BaseRepository
}

func NewRecommendationRepositorySimple(db *sql.DB) *RecommendationRepositorySimple {
	return &RecommendationRepositorySimple{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *RecommendationRepositorySimple) GetLatest(limit int) ([]*model.Recommendation, error) {
	latestRunAt, err := r.GetLatestRunAt()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest run_at: %w", err)
	}

	qb := NewQueryBuilder().
		Select("id", "ticker", "score", "explanation", "run_at", "rank").
		From("recommendations").
		Where("run_at = ?", latestRunAt).
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

func (r *RecommendationRepositorySimple) BulkCreate(recommendations []*model.Recommendation) error {
	if len(recommendations) == 0 {
		return nil
	}

	return r.ExecuteTransaction(func(tx *sql.Tx) error {
		query := `
			INSERT INTO recommendations (
				id, ticker, score, explanation, run_at, rank
			) VALUES ($1, $2, $3, $4, $5, $6)
		`

		stmt, err := tx.Prepare(query)
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %w", err)
		}
		defer stmt.Close()

		for _, recommendation := range recommendations {
			_, err := stmt.Exec(
				recommendation.ID, recommendation.Ticker, recommendation.Score,
				recommendation.Explanation, recommendation.RunAt, recommendation.Rank,
			)
			if err != nil {
				return fmt.Errorf("failed to insert recommendation %s: %w", recommendation.ID, err)
			}
		}

		return nil
	})
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

func (r *RecommendationRepositorySimple) ShouldCalculateToday() (bool, error) {
	latestRunAt, err := r.GetLatestRunAt()
	if err != nil {
		return true, nil
	}

	if latestRunAt == nil {
		return true, nil
	}

	today := time.Now().Format("2006-01-02")
	latestRunDate := latestRunAt.Format("2006-01-02")

	return today != latestRunDate, nil
}
