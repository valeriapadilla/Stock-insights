package response

import "time"

type HealthResponse struct {
    Status    string    `json:"status"`
    Timestamp time.Time `json:"timestamp"`
}

type AdminHealthResponse struct {
    Status             string    `json:"status"`
    Timestamp          time.Time `json:"timestamp"`
    DatabaseConnections int       `json:"database_connections,omitempty"`
    MemoryUsage        string    `json:"memory_usage,omitempty"`
    CacheHitRate       string    `json:"cache_hit_rate,omitempty"`
    ActiveRequests     int       `json:"active_requests,omitempty"`
    Uptime            string    `json:"uptime,omitempty"`
    Version           string    `json:"version"`
}