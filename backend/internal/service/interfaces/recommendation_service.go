package interfaces

import (
	"time"

	"github.com/valeriapadilla/stock-insights/internal/model"
)

type RecommendationService interface {
	GetLatestRecommendations(limit int) ([]*model.Recommendation, error)

	CreateRecommendations(recommendations []*model.Recommendation) error

	ShouldCalculateToday() (bool, error)

	GetLatestRunAt() (*time.Time, error)

	ValidateRecommendations(recommendations []*model.Recommendation) error
}
