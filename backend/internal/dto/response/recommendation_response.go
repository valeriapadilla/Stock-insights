package response

import (
    "time"
    "github.com/valeriapadilla/stock-insights/internal/model"
)

type RecommendationListResponse struct {
    Recommendations []model.Recommendation `json:"recommendations"`
    GeneratedAt     time.Time              `json:"generated_at"`
}