package repository

import (
	"database/sql"
	"time"

	"github.com/valeriapadilla/stock-insights/internal/model"
	"github.com/valeriapadilla/stock-insights/internal/repository/interfaces"
	"github.com/valeriapadilla/stock-insights/internal/validator"
)

type RecommendationRepository struct {
	*BaseRepository
	validator *validator.RecommendationValidator
}

var _ interfaces.RecommendationRepository = (*RecommendationRepository)(nil)

func NewRecommendationRepository(db *sql.DB) *RecommendationRepository {
	return &RecommendationRepository{
		BaseRepository: NewBaseRepository(db),
		validator:      validator.NewRecommendationValidator(),
	}
}

func (r *RecommendationRepository) CreateRecommendation(recommendation *model.Recommendation) error {
	query := `
		INSERT INTO recommendations (id, ticker, score, explanation, run_at, rank)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.GetDB().Exec(query,
		recommendation.ID,
		recommendation.Ticker,
		recommendation.Score,
		recommendation.Explanation,
		recommendation.RunAt,
		recommendation.Rank,
	)

	return err
}

func (r *RecommendationRepository) GetLatest(limit int) ([]*model.Recommendation, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT id, ticker, score, explanation, run_at, rank
		FROM recommendations
		WHERE run_at = (SELECT MAX(run_at) FROM recommendations)
		ORDER BY rank ASC
		LIMIT $1
	`

	rows, err := r.GetDB().Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recommendations []*model.Recommendation
	for rows.Next() {
		var rec model.Recommendation
		err := rows.Scan(
			&rec.ID,
			&rec.Ticker,
			&rec.Score,
			&rec.Explanation,
			&rec.RunAt,
			&rec.Rank,
		)
		if err != nil {
			return nil, err
		}
		recommendations = append(recommendations, &rec)
	}

	return recommendations, nil
}

func (r *RecommendationRepository) GetRecommendationsByDate(date time.Time, limit int) ([]*model.Recommendation, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT id, ticker, score, explanation, run_at, rank
		FROM recommendations
		WHERE DATE(run_at) = DATE($1)
		ORDER BY rank ASC
		LIMIT $2
	`

	rows, err := r.GetDB().Query(query, date, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recommendations []*model.Recommendation
	for rows.Next() {
		var rec model.Recommendation
		err := rows.Scan(
			&rec.ID,
			&rec.Ticker,
			&rec.Score,
			&rec.Explanation,
			&rec.RunAt,
			&rec.Rank,
		)
		if err != nil {
			return nil, err
		}
		recommendations = append(recommendations, &rec)
	}

	return recommendations, nil
}

func (r *RecommendationRepository) GetLatestRunAt() (*time.Time, error) {
	query := `SELECT MAX(run_at) FROM recommendations`

	var lastRun time.Time
	err := r.GetDB().QueryRow(query).Scan(&lastRun)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &lastRun, nil
}

func (r *RecommendationRepository) DeleteOldRecommendations(maxAge time.Duration) error {
	query := `DELETE FROM recommendations WHERE run_at < $1`

	cutoff := time.Now().Add(-maxAge)
	_, err := r.GetDB().Exec(query, cutoff)
	return err
}

func (r *RecommendationRepository) GetRecommendationCount() (int, error) {
	query := `SELECT COUNT(*) FROM recommendations`

	var count int
	err := r.GetDB().QueryRow(query).Scan(&count)
	return count, err
}
