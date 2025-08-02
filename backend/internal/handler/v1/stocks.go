package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/errors"
	"github.com/valeriapadilla/stock-insights/internal/service"
)

type StocksHandler struct {
	stockService *service.StockService
	logger       *logrus.Logger
}

func NewStocksHandler(stockService *service.StockService, logger *logrus.Logger) *StocksHandler {
	return &StocksHandler{
		stockService: stockService,
		logger:       logger,
	}
}

func (h *StocksHandler) ListStocks(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")
	sort := c.DefaultQuery("sort", "time")
	order := c.DefaultQuery("order", "desc")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	stocks, total, err := h.stockService.ListStocks(limit, offset, sort, order)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list stocks")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": "Failed to retrieve stocks",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stocks": stocks,
		"pagination": gin.H{
			"total":    total,
			"limit":    limit,
			"offset":   offset,
			"has_next": offset+limit < total,
		},
	})
}

func (h *StocksHandler) GetStock(c *gin.Context) {
	ticket := c.Param("ticket")
	if ticket == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": "Ticket parameter is required",
		})
		return
	}

	stock, err := h.stockService.GetStock(ticket)
	if err != nil {
		if errors.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not found",
				"message": "Stock not found",
			})
			return
		}

		h.logger.WithError(err).Error("Failed to get stock")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": "Failed to retrieve stock",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stock": stock,
	})
}

func (h *StocksHandler) SearchStocks(c *gin.Context) {
	ticket := c.Query("ticket")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var minPrice, maxPrice *float64
	if minPriceStr != "" {
		if val, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			minPrice = &val
		}
	}

	if maxPriceStr != "" {
		if val, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			maxPrice = &val
		}
	}

	stocks, total, err := h.stockService.SearchStocks(service.StockSearchParams{
		Ticket:   ticket,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		Limit:    limit,
		Offset:   offset,
	})

	if err != nil {
		h.logger.WithError(err).Error("Failed to search stocks")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": "Failed to search stocks",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stocks": stocks,
		"pagination": gin.H{
			"total":    total,
			"limit":    limit,
			"offset":   offset,
			"has_next": offset+limit < total,
		},
		"filters_applied": gin.H{
			"ticket":    ticket,
			"date_from": dateFrom,
			"date_to":   dateTo,
			"min_price": minPrice,
			"max_price": maxPrice,
		},
	})
}
