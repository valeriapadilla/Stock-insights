package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/app"
)

func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger := app.GetLogger().WithFields(logrus.Fields{
			"method":     param.Method,
			"path":       param.Path,
			"status":     param.StatusCode,
			"latency":    param.Latency.String(),
			"client_ip":  param.ClientIP,
			"user_agent": param.Request.UserAgent(),
			"timestamp":  param.TimeStamp.Format("2006-01-02T15:04:05.000Z"),
		})

		if param.StatusCode >= 400 {
			logger.Error("HTTP Request failed")
		} else {
			logger.Info("HTTP Request completed")
		}

		return ""
	})
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)

		logger := app.GetLogger().WithField("request_id", requestID)
		c.Set("logger", logger)

		c.Next()
	}
}

func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
