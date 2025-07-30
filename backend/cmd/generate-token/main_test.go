package main

import (
	"testing"
	"encoding/base64"
	"time"
	"crypto/rand"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/golang-jwt/jwt/v5"
)

// TestTokenGeneration verifica que se genere un token válido
func TestTokenGeneration(t *testing.T) {
	// Simular la función main
	secret := generateSecret()
	token := generateToken(secret)
	
	// Verificar que el secreto tiene 32 bytes (256 bits)
	secretBytes, err := base64.StdEncoding.DecodeString(secret)
	require.NoError(t, err)
	assert.Equal(t, 32, len(secretBytes))
	
	// Verificar que el token no está vacío
	assert.NotEmpty(t, token)
	assert.True(t, len(token) > 100) // JWT tokens son largos
}

// TestTokenValidation verifica que el token generado sea válido
func TestTokenValidation(t *testing.T) {
	secret := generateSecret()
	token := generateToken(secret)
	
	// Decodificar secreto
	secretBytes, err := base64.StdEncoding.DecodeString(secret)
	require.NoError(t, err)
	
	// Parsear y validar token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secretBytes, nil
	})
	
	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)
	
	// Verificar claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		assert.Equal(t, "admin", claims["role"])
		assert.Equal(t, "stock-insights-backend", claims["iss"])
		assert.Equal(t, "admin", claims["sub"])
	}
}

// Funciones helper para testing
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