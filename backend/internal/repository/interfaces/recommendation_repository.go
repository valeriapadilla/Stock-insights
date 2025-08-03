package interfaces

import (
	"database/sql"
	"time"

	"github.com/valeriapadilla/stock-insights/internal/model"
)

type RecommendationRepository interface {
	CreateRecommendation(recommendation *model.Recommendation) error
	GetLatest(limit int) ([]*model.Recommendation, error)
	GetLatestRunAt() (*time.Time, error)
	GetDB() *sql.DB
}
