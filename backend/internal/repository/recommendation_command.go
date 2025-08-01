package repository

import (
	"database/sql"
	"fmt"

	"github.com/valeriapadilla/stock-insights/internal/model"
	"github.com/valeriapadilla/stock-insights/internal/repository/interfaces"
	"github.com/valeriapadilla/stock-insights/internal/validator"
)

type RecommendationCommandImpl struct {
	*BaseRepository
	validator *validator.RecommendationValidator
}

var _ interfaces.RecommendationCommand = (*RecommendationCommandImpl)(nil)

func NewRecommendationCommand(db *sql.DB) *RecommendationCommandImpl {
	return &RecommendationCommandImpl{
		BaseRepository: NewBaseRepository(db),
		validator:      validator.NewRecommendationValidator(),
	}
}

func (c *RecommendationCommandImpl) BulkCreate(recommendations []*model.Recommendation) error {
	if len(recommendations) == 0 {
		return nil
	}

	tx, err := c.GetDB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO recommendations (ticker, score, explanation, rank, run_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, recommendation := range recommendations {
		if err := c.validateRecommendation(recommendation); err != nil {
			return fmt.Errorf("recommendation validation failed for %s: %w", recommendation.Ticker, err)
		}

		_, err := stmt.Exec(
			recommendation.Ticker, recommendation.Score,
			recommendation.Explanation, recommendation.Rank, recommendation.RunAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert recommendation %s: %w", recommendation.Ticker, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (c *RecommendationCommandImpl) validateRecommendation(recommendation *model.Recommendation) error {
	if recommendation == nil {
		return fmt.Errorf("recommendation cannot be nil")
	}

	if recommendation.Ticker == "" {
		return fmt.Errorf("ticker is required")
	}
	if err := c.validator.Validate(recommendation); err != nil {
		return fmt.Errorf("recommendation validation failed: %w", err)
	}

	return nil
}
