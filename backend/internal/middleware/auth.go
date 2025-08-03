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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "You are not allowed to access this resource. Authentication required.",
			})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid authentication format. Use Bearer token.",
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "JWT secret not configured",
			})
			c.Abort()
			return
		}

		secretBytes, err := base64.StdEncoding.DecodeString(secret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid JWT secret format",
			})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretBytes, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid or expired authentication token.",
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid or expired authentication token.",
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if role, exists := claims["role"]; !exists || role != "admin" {
				c.JSON(http.StatusForbidden, gin.H{
					"error":   "Forbidden",
					"message": "You are not allowed to access this resource. Admin privileges required.",
				})
				c.Abort()
				return
			}

			c.Set("user_role", claims["role"])
			c.Set("user_subject", claims["sub"])
		}

		c.Next()
	}
}
