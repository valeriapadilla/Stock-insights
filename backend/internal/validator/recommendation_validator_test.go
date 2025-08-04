package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valeriapadilla/stock-insights/internal/model"
)

func TestRecommendationValidator_Validate(t *testing.T) {
	validator := NewRecommendationValidator()

	tests := []struct {
		name    string
		rec     *model.Recommendation
		wantErr bool
	}{
		{
			name: "valid recommendation",
			rec: &model.Recommendation{
				ID:          "test-id",
				Ticker:      "AAPL",
				Score:       85.5,
				Explanation: "Test explanation",
				Rank:        1,
			},
			wantErr: false,
		},
		{
			name:    "nil recommendation",
			rec:     nil,
			wantErr: true,
		},
		{
			name: "empty ticker",
			rec: &model.Recommendation{
				ID:          "test-id",
				Ticker:      "",
				Score:       85.5,
				Explanation: "Test explanation",
				Rank:        1,
			},
			wantErr: true,
		},
		{
			name: "ticker too long",
			rec: &model.Recommendation{
				ID:          "test-id",
				Ticker:      "VERYLONGTICKER",
				Score:       85.5,
				Explanation: "Test explanation",
				Rank:        1,
			},
			wantErr: true,
		},
		{
			name: "score too high",
			rec: &model.Recommendation{
				ID:          "test-id",
				Ticker:      "AAPL",
				Score:       150.0,
				Explanation: "Test explanation",
				Rank:        1,
			},
			wantErr: true,
		},
		{
			name: "score negative",
			rec: &model.Recommendation{
				ID:          "test-id",
				Ticker:      "AAPL",
				Score:       -10.0,
				Explanation: "Test explanation",
				Rank:        1,
			},
			wantErr: true,
		},
		{
			name: "explanation too long",
			rec: &model.Recommendation{
				ID:          "test-id",
				Ticker:      "AAPL",
				Score:       85.5,
				Explanation: "This is a very long explanation that exceeds the maximum allowed length of 500 characters. This is a very long explanation that exceeds the maximum allowed length of 500 characters. This is a very long explanation that exceeds the maximum allowed length of 500 characters. This is a very long explanation that exceeds the maximum allowed length of 500 characters. This is a very long explanation that exceeds the maximum allowed length of 500 characters. This is a very long explanation that exceeds the maximum allowed length of 500 characters. This is a very long explanation that exceeds the maximum allowed length of 500 characters. This is a very long explanation that exceeds the maximum allowed length of 500 characters. This is a very long explanation that exceeds the maximum allowed length of 500 characters. This is a very long explanation that exceeds the maximum allowed length of 500 characters.",
				Rank:        1,
			},
			wantErr: true,
		},
		{
			name: "rank zero",
			rec: &model.Recommendation{
				ID:          "test-id",
				Ticker:      "AAPL",
				Score:       85.5,
				Explanation: "Test explanation",
				Rank:        0,
			},
			wantErr: true,
		},
		{
			name: "rank negative",
			rec: &model.Recommendation{
				ID:          "test-id",
				Ticker:      "AAPL",
				Score:       85.5,
				Explanation: "Test explanation",
				Rank:        -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.rec)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRecommendationValidator_ValidateBulk(t *testing.T) {
	validator := NewRecommendationValidator()

	tests := []struct {
		name    string
		recs    []*model.Recommendation
		wantErr bool
	}{
		{
			name: "valid bulk",
			recs: []*model.Recommendation{
				{
					ID:          "test-id-1",
					Ticker:      "AAPL",
					Score:       85.5,
					Explanation: "Test explanation 1",
					Rank:        1,
				},
				{
					ID:          "test-id-2",
					Ticker:      "GOOGL",
					Score:       90.0,
					Explanation: "Test explanation 2",
					Rank:        2,
				},
			},
			wantErr: false,
		},
		{
			name:    "empty slice",
			recs:    []*model.Recommendation{},
			wantErr: false,
		},
		{
			name: "invalid recommendation in bulk",
			recs: []*model.Recommendation{
				{
					ID:          "test-id-1",
					Ticker:      "AAPL",
					Score:       85.5,
					Explanation: "Test explanation 1",
					Rank:        1,
				},
				{
					ID:          "test-id-2",
					Ticker:      "", // Invalid
					Score:       90.0,
					Explanation: "Test explanation 2",
					Rank:        2,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateBulk(tt.recs)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRecommendationValidator_ValidateRecommendationParams(t *testing.T) {
	validator := NewRecommendationValidator()

	tests := []struct {
		name   string
		params RecommendationParams
		want   RecommendationParams
	}{
		{
			name: "valid params",
			params: RecommendationParams{
				DaysBack:   7,
				MaxResults: 30,
				MinScore:   80,
			},
			want: RecommendationParams{
				DaysBack:   7,
				MaxResults: 30,
				MinScore:   80,
			},
		},
		{
			name: "zero values use defaults",
			params: RecommendationParams{
				DaysBack:   0,
				MaxResults: 0,
				MinScore:   0,
			},
			want: RecommendationParams{
				DaysBack:   7,
				MaxResults: 0,
				MinScore:   0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateRecommendationParams(tt.params)
			assert.Equal(t, tt.want, result)
		})
	}
}
