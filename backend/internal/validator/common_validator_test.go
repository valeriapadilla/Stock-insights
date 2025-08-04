package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommonValidator_SanitizeString(t *testing.T) {
	validator := NewCommonValidator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal string",
			input:    "test string",
			expected: "test string",
		},
		{
			name:     "string with spaces",
			input:    "  test string  ",
			expected: "test string",
		},
		{
			name:     "string too long",
			input:    "this is a very long string that exceeds the maximum allowed length of 100 characters and should be truncated",
			expected: "this is a very long string that exceeds the maximum allowed length of 100 characters and should be t",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.SanitizeString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCommonValidator_ValidateFilter(t *testing.T) {
	validator := NewCommonValidator()

	tests := []struct {
		name        string
		key         string
		value       string
		expectedErr bool
	}{
		{
			name:        "valid filter",
			key:         "ticker",
			value:       "AAPL",
			expectedErr: false,
		},
		{
			name:        "empty value",
			key:         "ticker",
			value:       "",
			expectedErr: false,
		},
		{
			name:        "value with spaces",
			key:         "ticker",
			value:       "  AAPL  ",
			expectedErr: true,
		},
		{
			name:        "value too long",
			key:         "ticker",
			value:       "this is a very long string that exceeds the maximum allowed length of 100 characters and should be truncated",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateFilter(tt.key, tt.value)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCommonValidator_ValidateFilters(t *testing.T) {
	validator := NewCommonValidator()

	tests := []struct {
		name        string
		filters     map[string]string
		expectedErr bool
	}{
		{
			name: "valid filters",
			filters: map[string]string{
				"ticker": "AAPL",
				"date":   "2025-01-01",
			},
			expectedErr: false,
		},
		{
			name:        "empty filters",
			filters:     map[string]string{},
			expectedErr: false,
		},
		{
			name: "invalid filter",
			filters: map[string]string{
				"ticker": "  AAPL  ",
			},
			expectedErr: true,
		},
		{
			name: "mixed valid and invalid filters",
			filters: map[string]string{
				"ticker": "AAPL",
				"date":   "  invalid  ",
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateFilters(tt.filters)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
