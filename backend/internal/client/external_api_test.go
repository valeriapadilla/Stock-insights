package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewExternalAPIClient(t *testing.T) {
	logger := logrus.New()
	config := ExternalAPIConfig{
		BaseURL:    "https://api.karenai.click/swechallenge/list",
		APIKey:     "Bearer test-key",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
	}

	client := NewExternalAPIClient(config, logger)

	assert.NotNil(t, client)
	assert.Equal(t, config.BaseURL, client.baseURL)
	assert.Equal(t, config.APIKey, client.apiKey)
	assert.NotNil(t, client.httpClient)
	assert.Equal(t, config.Timeout, client.httpClient.Timeout)
}

func TestExternalAPIClient_GetAllStocks_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/", r.URL.String())
		assert.Equal(t, "test-key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		response := `{
			"items": [
				{
					"ticker": "AAPL",
					"company": "Apple Inc.",
					"target_from": "$150.00",
					"target_to": "$160.00",
					"rating_from": "Buy",
					"rating_to": "Buy",
					"action": "initiated by",
					"brokerage": "Goldman Sachs",
					"time": "2024-01-15T10:30:00Z",
					"created_at": "2024-01-15T10:30:00Z",
					"updated_at": "2024-01-15T10:30:00Z"
				},
				{
					"ticker": "GOOGL",
					"company": "Alphabet Inc.",
					"target_from": "$2800.00",
					"target_to": "$2900.00",
					"rating_from": "Buy",
					"rating_to": "Buy",
					"action": "initiated by",
					"brokerage": "Morgan Stanley",
					"time": "2024-01-15T10:30:00Z",
					"created_at": "2024-01-15T10:30:00Z",
					"updated_at": "2024-01-15T10:30:00Z"
				}
			]
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	logger := logrus.New()
	config := ExternalAPIConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 30 * time.Second,
	}

	client := NewExternalAPIClient(config, logger)

	ctx := context.Background()
	stocks, err := client.GetAllStocks(ctx)

	require.NoError(t, err)
	assert.Len(t, stocks, 2)

	stock1 := stocks[0]
	assert.Equal(t, "AAPL", stock1.Ticker)
	assert.Equal(t, "Apple Inc.", stock1.Company)
	assert.Equal(t, "$150.00", stock1.TargetFrom)
	assert.Equal(t, "$160.00", stock1.TargetTo)
	assert.Equal(t, "Buy", stock1.RatingFrom)
	assert.Equal(t, "Buy", stock1.RatingTo)
	assert.Equal(t, "initiated by", stock1.Action)
	assert.Equal(t, "Goldman Sachs", stock1.Brokerage)

	stock2 := stocks[1]
	assert.Equal(t, "GOOGL", stock2.Ticker)
	assert.Equal(t, "Alphabet Inc.", stock2.Company)
}

func TestExternalAPIClient_GetAllStocks_WithPagination(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		var response string

		switch callCount {
		case 1:
			assert.Equal(t, "test-key", r.Header.Get("Authorization"))
			response = `{
				"items": [
					{"ticker": "AAPL", "company": "Apple Inc.", "target_from": "$150", "target_to": "$160", "rating_from": "Buy", "rating_to": "Buy", "action": "initiated", "brokerage": "Goldman", "time": "2024-01-15T10:30:00Z", "created_at": "2024-01-15T10:30:00Z", "updated_at": "2024-01-15T10:30:00Z"},
					{"ticker": "GOOGL", "company": "Alphabet Inc.", "target_from": "$2800", "target_to": "$2900", "rating_from": "Buy", "rating_to": "Buy", "action": "initiated", "brokerage": "Morgan Stanley", "time": "2024-01-15T10:30:00Z", "created_at": "2024-01-15T10:30:00Z", "updated_at": "2024-01-15T10:30:00Z"}
				],
				"next_page": "GOOGL"
			}`
		case 2:
			assert.Equal(t, "/?next_page=GOOGL", r.URL.String())
			assert.Equal(t, "test-key", r.Header.Get("Authorization"))
			response = `{
				"items": [
					{"ticker": "MSFT", "company": "Microsoft Corp.", "target_from": "$300", "target_to": "$320", "rating_from": "Buy", "rating_to": "Buy", "action": "initiated", "brokerage": "JP Morgan", "time": "2024-01-15T10:30:00Z", "created_at": "2024-01-15T10:30:00Z", "updated_at": "2024-01-15T10:30:00Z"}
				]
			}`
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	logger := logrus.New()
	config := ExternalAPIConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 30 * time.Second,
	}

	client := NewExternalAPIClient(config, logger)

	ctx := context.Background()
	stocks, err := client.GetAllStocks(ctx)

	require.NoError(t, err)
	assert.Len(t, stocks, 3)
	assert.Equal(t, "AAPL", stocks[0].Ticker)
	assert.Equal(t, "GOOGL", stocks[1].Ticker)
	assert.Equal(t, "MSFT", stocks[2].Ticker)
	assert.Equal(t, 2, callCount)
}

func TestExternalAPIClient_GetAllStocks_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))
	}))
	defer server.Close()

	logger := logrus.New()
	config := ExternalAPIConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 30 * time.Second,
	}

	client := NewExternalAPIClient(config, logger)

	ctx := context.Background()
	stocks, err := client.GetAllStocks(ctx)

	assert.Error(t, err)
	assert.Nil(t, stocks)
	assert.Contains(t, err.Error(), "API returned status 500")
}

func TestExternalAPIClient_GetAllStocks_NetworkError(t *testing.T) {
	logger := logrus.New()
	config := ExternalAPIConfig{
		BaseURL: "http://invalid-url-that-does-not-exist.com",
		APIKey:  "test-key",
		Timeout: 1 * time.Second,
	}

	client := NewExternalAPIClient(config, logger)

	ctx := context.Background()
	stocks, err := client.GetAllStocks(ctx)

	assert.Error(t, err)
	assert.Nil(t, stocks)
	assert.True(t, strings.Contains(err.Error(), "failed to make HTTP request") ||
		strings.Contains(err.Error(), "failed to decode response") ||
		strings.Contains(err.Error(), "no such host"))
}

func TestExternalAPIClient_GetAllStocks_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"invalid": json`))
	}))
	defer server.Close()

	logger := logrus.New()
	config := ExternalAPIConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 30 * time.Second,
	}

	client := NewExternalAPIClient(config, logger)

	ctx := context.Background()
	stocks, err := client.GetAllStocks(ctx)

	assert.Error(t, err)
	assert.Nil(t, stocks)
	assert.Contains(t, err.Error(), "failed to decode response")
}

func TestExternalAPIClient_HealthCheck_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/", r.URL.String())
		assert.Equal(t, "test-key", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	}))
	defer server.Close()

	logger := logrus.New()
	config := ExternalAPIConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 30 * time.Second,
	}

	client := NewExternalAPIClient(config, logger)

	ctx := context.Background()
	err := client.HealthCheck(ctx)

	assert.NoError(t, err)
}

func TestExternalAPIClient_HealthCheck_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status": "unhealthy"}`))
	}))
	defer server.Close()

	logger := logrus.New()
	config := ExternalAPIConfig{
		BaseURL: server.URL,
		APIKey:  "test-key",
		Timeout: 30 * time.Second,
	}

	client := NewExternalAPIClient(config, logger)

	ctx := context.Background()
	err := client.HealthCheck(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "health check returned status 503")
}

func TestExternalAPIClient_WithoutAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get("Authorization"))

		response := `{
			"items": []
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	logger := logrus.New()
	config := ExternalAPIConfig{
		BaseURL: server.URL,
		APIKey:  "",
		Timeout: 30 * time.Second,
	}

	client := NewExternalAPIClient(config, logger)

	ctx := context.Background()
	stocks, err := client.GetAllStocks(ctx)

	require.NoError(t, err)
	assert.Len(t, stocks, 0)
}
