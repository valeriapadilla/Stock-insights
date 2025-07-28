package response

import "time"

type AdminResponse struct {
    Status    string    `json:"status"`
    Message   string    `json:"message"`
    Timestamp time.Time `json:"timestamp"`
}

type StatsResponse struct {
    TotalStocks        int    `json:"total_stocks"`
    LastIngestion      string `json:"last_ingestion"`
    LastRecommendations string `json:"last_recommendations"`
    CacheStats         struct {
        HitRate    string `json:"hit_rate"`
        TotalItems int    `json:"total_items"`
    } `json:"cache_stats"`
    APIStats struct {
        RequestsToday    int    `json:"requests_today"`
        AvgResponseTime  string `json:"avg_response_time"`
    } `json:"api_stats"`
}