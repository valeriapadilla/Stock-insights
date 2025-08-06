package utils

import (
	"math"
	"strconv"
	"strings"
)

func ParsePrice(priceStr string) float64 {
	if priceStr == "" {
		return 0.0
	}

	cleanPrice := strings.TrimSpace(strings.ReplaceAll(priceStr, "$", ""))

	price, err := strconv.ParseFloat(cleanPrice, 64)
	if err != nil {
		return 0.0
	}

	return price
}

func CalculateChangePercentage(fromPrice, toPrice float64) float64 {
	if fromPrice <= 0 {
		return 0.0
	}

	change := ((toPrice - fromPrice) / fromPrice) * 100
	return math.Round(change*100) / 100
}
