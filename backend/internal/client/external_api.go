package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/valeriapadilla/stock-insights/internal/errors"
	"github.com/valeriapadilla/stock-insights/internal/model"
)

type ExternalAPIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	logger     *logrus.Logger
}

type ExternalAPIConfig struct {
	BaseURL    string
	APIKey     string
	Timeout    time.Duration
	MaxRetries int
}

func NewExternalAPIClient(config ExternalAPIConfig, logger *logrus.Logger) *ExternalAPIClient {
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &ExternalAPIClient{
		baseURL:    config.BaseURL,
		apiKey:     config.APIKey,
		httpClient: httpClient,
		logger:     logger,
	}
}

func (c *ExternalAPIClient) GetAllStocks(ctx context.Context) ([]model.Stock, error) {
	var allStocks []model.Stock
	nextPage := ""
	pageCount := 0

	for {
		url := c.baseURL
		if nextPage != "" {
			url = fmt.Sprintf("%s?next_page=%s", c.baseURL, nextPage)
		}

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, errors.NewInternalError("failed to create request", err)
		}

		if c.apiKey != "" {
			req.Header.Set("Authorization", c.apiKey)
		}
		req.Header.Set("Content-Type", "application/json")

		c.logger.WithFields(logrus.Fields{
			"url":       url,
			"method":    "GET",
			"page":      pageCount + 1,
			"next_page": nextPage,
		}).Debug("Making external API request")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			c.logger.WithError(err).Error("HTTP request failed")
			return nil, errors.NewExternalError("failed to make HTTP request", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errorMsg := fmt.Sprintf("API returned status %d: %s", resp.StatusCode, string(body))

			c.logger.WithFields(logrus.Fields{
				"status_code":   resp.StatusCode,
				"response_body": string(body),
				"url":           url,
			}).Error("External API error")

			return nil, errors.NewExternalError(errorMsg, nil)
		}

		var apiResponse model.ExternalAPIResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			c.logger.WithError(err).Error("Failed to decode API response")
			return nil, errors.NewInternalError("failed to decode response", err)
		}

		allStocks = append(allStocks, apiResponse.Items...)
		pageCount++

		c.logger.WithFields(logrus.Fields{
			"page":          pageCount,
			"items_in_page": len(apiResponse.Items),
			"total_items":   len(allStocks),
			"has_next_page": apiResponse.NextPage != "",
		}).Info("Retrieved page from external API")

		if apiResponse.NextPage == "" {
			break
		}

		nextPage = apiResponse.NextPage
	}

	c.logger.WithFields(logrus.Fields{
		"total_stocks": len(allStocks),
		"total_pages":  pageCount,
	}).Info("Successfully retrieved all stocks from external API")

	return allStocks, nil
}

func (c *ExternalAPIClient) HealthCheck(ctx context.Context) error {
	url := c.baseURL

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return errors.NewInternalError("failed to create health check request", err)
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.WithError(err).Error("Health check HTTP request failed")
		return errors.NewExternalError("health check failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		errorMsg := fmt.Sprintf("health check returned status %d: %s", resp.StatusCode, string(body))

		c.logger.WithFields(logrus.Fields{
			"status_code":   resp.StatusCode,
			"response_body": string(body),
			"url":           url,
		}).Error("External API health check failed")

		return errors.NewExternalError(errorMsg, nil)
	}

	c.logger.Info("External API health check passed")
	return nil
}
