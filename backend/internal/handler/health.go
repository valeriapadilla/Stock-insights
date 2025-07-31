package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/valeriapadilla/stock-insights/internal/dto/response"
)

func HealthCheck(c *gin.Context) {
	healthResp := response.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, healthResp)
}

func HealthCheckHead(c *gin.Context) {
	c.Status(http.StatusOK)
}
