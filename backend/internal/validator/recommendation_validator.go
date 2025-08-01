package validator

import (
	"fmt"

	"github.com/valeriapadilla/stock-insights/internal/model"
)

type RecommendationValidator struct{}

func NewRecommendationValidator() *RecommendationValidator {
	return &RecommendationValidator{}
}

func (v *RecommendationValidator) Validate(rec *model.Recommendation) error {
	if rec == nil {
		return fmt.Errorf("recommendation cannot be nil")
	}

	if err := v.validateID(rec.ID); err != nil {
		return fmt.Errorf("ID: %w", err)
	}

	if err := v.validateTicker(rec.Ticker); err != nil {
		return fmt.Errorf("Ticker: %w", err)
	}

	if err := v.validateScore(rec.Score); err != nil {
		return fmt.Errorf("Score: %w", err)
	}

	if err := v.validateExplanation(rec.Explanation); err != nil {
		return fmt.Errorf("Explanation: %w", err)
	}

	if err := v.validateRank(rec.Rank); err != nil {
		return fmt.Errorf("Rank: %w", err)
	}

	return nil
}

func (v *RecommendationValidator) ValidateBulk(recommendations []*model.Recommendation) error {
	if len(recommendations) == 0 {
		return nil
	}

	for i, rec := range recommendations {
		if err := v.Validate(rec); err != nil {
			return fmt.Errorf("recommendation[%d]: %w", i, err)
		}
	}

	return nil
}

func (v *RecommendationValidator) validateID(id string) error {
	// Allow empty ID as database will generate UUID
	if id != "" && len(id) > 50 {
		return fmt.Errorf("must be 50 characters or less, got %d", len(id))
	}
	return nil
}

func (v *RecommendationValidator) validateTicker(ticker string) error {
	if ticker == "" {
		return fmt.Errorf("cannot be empty")
	}
	if len(ticker) > 10 {
		return fmt.Errorf("must be 10 characters or less, got %d", len(ticker))
	}
	return nil
}

func (v *RecommendationValidator) validateScore(score float64) error {
	if score < 0 {
		return fmt.Errorf("must be 0 or greater, got %.2f", score)
	}
	if score > 100 {
		return fmt.Errorf("must be 100 or less, got %.2f", score)
	}
	return nil
}

func (v *RecommendationValidator) validateExplanation(explanation string) error {
	if len(explanation) > 500 {
		return fmt.Errorf("must be 500 characters or less, got %d", len(explanation))
	}
	return nil
}

func (v *RecommendationValidator) validateRank(rank int) error {
	if rank <= 0 {
		return fmt.Errorf("must be greater than 0, got %d", rank)
	}
	return nil
}
