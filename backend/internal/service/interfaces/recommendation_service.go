package interfaces

import (
	"github.com/valeriapadilla/stock-insights/internal/model"
	"github.com/valeriapadilla/stock-insights/internal/validator"
)

type RecommendationServiceInterface interface {
	CalculateRecommendations(params validator.RecommendationParams) ([]*model.Recommendation, error)
	GetLatestRecommendations(limit int) ([]*model.Recommendation, error)
	SaveRecommendations(recommendations []*model.Recommendation) error
}
