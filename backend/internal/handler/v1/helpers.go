package v1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/errors"
)

func parsePaginationParams(c *gin.Context) (limit, offset int, sort, order string) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")
	sort = c.DefaultQuery("sort", "time")
	order = c.DefaultQuery("order", "desc")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	offset, err = strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return limit, offset, sort, order
}

func parsePriceParams(c *gin.Context) (*float64, *float64) {
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")

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

	return minPrice, maxPrice
}

func parseFilterParams(c *gin.Context) (rating, sortBy, order string) {
	rating = strings.ToLower(c.Query("rating"))
	sortBy = c.DefaultQuery("sort_by", "time")
	order = c.DefaultQuery("order", "desc")

	validSortBy := map[string]bool{
		"ticker":         true,
		"change_percent": true,
		"time":           true,
	}

	if !validSortBy[sortBy] {
		sortBy = "time"
	}

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return rating, sortBy, order
}

func handleError(c *gin.Context, err error, operation string, logger *logrus.Logger) {
	logger.WithError(err).Errorf("Failed to %s", operation)

	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.Code, gin.H{
			"error":   appErr.Type,
			"message": appErr.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   "Internal server error",
		"message": "Failed to " + operation,
	})
}
