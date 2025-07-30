package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// TestAuthMiddlewareValidToken verifica middleware con token válido
func TestAuthMiddlewareValidToken(t *testing.T) {
	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)
	
	// Crear router con middleware
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// Generar token válido
	secret := generateTestSecret()
	token := generateTestToken(secret)
	
	// Configurar JWT_SECRET en environment
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")
	
	// Crear request con token válido
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	// Ejecutar request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Verificar respuesta exitosa
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestAuthMiddlewareInvalidToken verifica middleware con token inválido
func TestAuthMiddlewareInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// Configurar JWT_SECRET
	secret := generateTestSecret()
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")
	
	// Crear request con token inválido
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Verificar respuesta de error
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthMiddlewareNoToken verifica middleware sin token
func TestAuthMiddlewareNoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// Crear request sin token
	req, _ := http.NewRequest("GET", "/test", nil)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Verificar respuesta de error
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAuthMiddlewareWrongFormat verifica middleware con formato incorrecto
func TestAuthMiddlewareWrongFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	
	// Crear request con formato incorrecto
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat token")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Verificar respuesta de error
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Funciones helper para testing
func generateTestSecret() string {
	secret := make([]byte, 32)
	rand.Read(secret)
	return base64.StdEncoding.EncodeToString(secret)
}

func generateTestToken(secret string) string {
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