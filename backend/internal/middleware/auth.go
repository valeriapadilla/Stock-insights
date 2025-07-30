package middleware

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware verifica el token JWT en headers
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authorization header required",
				"message": "Use format: Bearer <token>",
			})
			c.Abort()
			return
		}

		// Verificar formato "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid authorization format",
				"message": "Use format: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Obtener secreto JWT del environment
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "JWT secret not configured",
			})
			c.Abort()
			return
		}

		// Decodificar secreto de base64
		secretBytes, err := base64.StdEncoding.DecodeString(secret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid JWT secret format",
			})
			c.Abort()
			return
		}

		// Verificar y parsear token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verificar algoritmo
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretBytes, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// Verificar que el token es válido
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		// Verificar claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Verificar rol admin
			if role, exists := claims["role"]; !exists || role != "admin" {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Admin role required",
				})
				c.Abort()
				return
			}

			// Agregar claims al contexto para uso posterior
			c.Set("user_role", claims["role"])
			c.Set("user_subject", claims["sub"])
		}

		c.Next()
	}
}