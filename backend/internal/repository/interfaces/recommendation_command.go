package interfaces

import (
	"github.com/valeriapadilla/stock-insights/internal/model"
)

type RecommendationCommand interface {
	BulkCreate(recommendations []*model.Recommendation) error
}
