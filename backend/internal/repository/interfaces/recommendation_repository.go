package interfaces

import (
	"database/sql"
	"time"

	"github.com/valeriapadilla/stock-insights/internal/model"
)

type RecommendationRepository interface {
	GetLatest(limit int) ([]*model.Recommendation, error)

	BulkCreate(recommendations []*model.Recommendation) error

	GetLatestRunAt() (*time.Time, error)

	GetDB() *sql.DB
}
