package response

import (
    "github.com/valeriapadilla/stock-insights/internal/model"
)

type StockListResponse struct {
    Stocks     []model.Stock `json:"stocks"`
    Pagination struct {
        Page       int `json:"page"`
        Limit      int `json:"limit"`
        Total      int `json:"total"`
        TotalPages int `json:"total_pages"`
    } `json:"pagination"`
}

type StockDetailResponse struct {
    Stock model.Stock `json:"stock"`
}