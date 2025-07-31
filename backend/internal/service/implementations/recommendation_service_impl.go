package implementations

import (
	"fmt"
	"time"

	"github.com/valeriapadilla/stock-insights/internal/model"
	repoInterfaces "github.com/valeriapadilla/stock-insights/internal/repository/interfaces"
	serviceInterfaces "github.com/valeriapadilla/stock-insights/internal/service/interfaces"
	"github.com/valeriapadilla/stock-insights/internal/validator"
)

type RecommendationServiceImpl struct {
	recommendationRepo repoInterfaces.RecommendationRepository
	validator          *validator.RecommendationValidator
}

var _ serviceInterfaces.RecommendationService = (*RecommendationServiceImpl)(nil)

func NewRecommendationService(recommendationRepo repoInterfaces.RecommendationRepository) *RecommendationServiceImpl {
	return &RecommendationServiceImpl{
		recommendationRepo: recommendationRepo,
		validator:          validator.NewRecommendationValidator(),
	}
}

func (s *RecommendationServiceImpl) GetLatestRecommendations(limit int) ([]*model.Recommendation, error) {
	limit = s.applyRecommendationLimit(limit)

	recommendations, err := s.recommendationRepo.GetLatest(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest recommendations: %w", err)
	}

	return recommendations, nil
}

func (s *RecommendationServiceImpl) CreateRecommendations(recommendations []*model.Recommendation) error {
	if err := s.ValidateRecommendations(recommendations); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := s.applyBulkCreationRules(recommendations); err != nil {
		return fmt.Errorf("business rules validation failed: %w", err)
	}

	if err := s.recommendationRepo.BulkCreate(recommendations); err != nil {
		return fmt.Errorf("failed to create recommendations: %w", err)
	}

	return nil
}

func (s *RecommendationServiceImpl) ShouldCalculateToday() (bool, error) {
	latestRunAt, err := s.recommendationRepo.GetLatestRunAt()
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

func (s *RecommendationServiceImpl) GetLatestRunAt() (*time.Time, error) {
	return s.recommendationRepo.GetLatestRunAt()
}

func (s *RecommendationServiceImpl) ValidateRecommendations(recommendations []*model.Recommendation) error {
	if len(recommendations) == 0 {
		return nil
	}

	if err := s.applyBulkValidationRules(recommendations); err != nil {
		return fmt.Errorf("business rules validation failed: %w", err)
	}

	return s.validator.ValidateBulk(recommendations)
}

func (s *RecommendationServiceImpl) applyRecommendationLimit(limit int) int {
	if limit <= 0 {
		return 10
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func (s *RecommendationServiceImpl) applyBulkCreationRules(recommendations []*model.Recommendation) error {
	if len(recommendations) > 1000 {
		return fmt.Errorf("cannot create more than 1000 recommendations at once")
	}

	if len(recommendations) > 1 {
		firstRunAt := recommendations[0].RunAt
		for i, rec := range recommendations {
			if !rec.RunAt.Equal(firstRunAt) {
				return fmt.Errorf("all recommendations must have the same run_at date, recommendation[%d] differs", i)
			}
		}
	}

	return nil
}

func (s *RecommendationServiceImpl) applyBulkValidationRules(recommendations []*model.Recommendation) error {
	tickerMap := make(map[string]bool)
	for _, rec := range recommendations {
		if tickerMap[rec.Ticker] {
			return fmt.Errorf("duplicate ticker '%s' found in recommendations", rec.Ticker)
		}
		tickerMap[rec.Ticker] = true
	}

	return nil
}
