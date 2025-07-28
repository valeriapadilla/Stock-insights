package response

import (
    "time"
    "stock-insights/backend/internal/model"
)

type RecommendationListResponse struct {
    Recommendations []model.Recommendation `json:"recommendations"`
    GeneratedAt     time.Time              `json:"generated_at"`
}