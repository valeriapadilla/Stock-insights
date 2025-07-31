package main

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenGeneration(t *testing.T) {
	secret := generateSecret()
	token := generateToken(secret)

	secretBytes, err := base64.StdEncoding.DecodeString(secret)
	require.NoError(t, err)
	assert.Equal(t, 32, len(secretBytes))

	assert.NotEmpty(t, token)
	assert.True(t, len(token) > 100)
}

func TestTokenValidation(t *testing.T) {
	secret := generateSecret()
	token := generateToken(secret)

	secretBytes, err := base64.StdEncoding.DecodeString(secret)
	require.NoError(t, err)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secretBytes, nil
	})

	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		assert.Equal(t, "admin", claims["role"])
		assert.Equal(t, "stock-insights-backend", claims["iss"])
		assert.Equal(t, "admin", claims["sub"])
	}
}

func generateSecret() string {
	secret := make([]byte, 32)
	rand.Read(secret)
	return base64.StdEncoding.EncodeToString(secret)
}

func generateToken(secret string) string {
	secretBytes, _ := base64.StdEncoding.DecodeString(secret)

	claims := jwt.MapClaims{
		"role": "admin",
		"exp":  time.Now().AddDate(0, 1, 0).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  "stock-insights-backend",
		"sub":  "admin",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(secretBytes)
	return tokenString
}
