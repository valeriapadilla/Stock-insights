package validator

import (
	"fmt"
	"strings"
)

type CommonValidator struct{}

func NewCommonValidator() *CommonValidator {
	return &CommonValidator{}
}

func (v *CommonValidator) SanitizeString(input string) string {
	input = strings.TrimSpace(input)
	if len(input) > 100 {
		input = input[:100]
	}
	return input
}

func (v *CommonValidator) ValidateFilter(key, value string) error {
	if value == "" {
		return nil
	}

	sanitized := v.SanitizeString(value)
	if sanitized != value {
		return fmt.Errorf("filter value '%s' contains invalid characters", key)
	}

	return nil
}

func (v *CommonValidator) ValidateFilters(filters map[string]string) error {
	for key, value := range filters {
		if err := v.ValidateFilter(key, value); err != nil {
			return fmt.Errorf("filter '%s': %w", key, err)
		}
	}

	return nil
}
