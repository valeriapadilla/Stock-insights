package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	fmt.Println("Setting up authentication...")

	jwtSecret := generateJWTSecret()
	fmt.Printf("Generated JWT_SECRET: %s\n", jwtSecret[:20]+"...")

	adminToken := generateAdminToken(jwtSecret)
	fmt.Printf("Generated admin token: %s\n", adminToken[:20]+"...")

	updateEnvFile(jwtSecret, adminToken)

	fmt.Println("Authentication setup completed!")
	fmt.Println("You can now use the admin token to access protected endpoints")
	fmt.Println("Example: curl -X POST http://localhost:8080/api/v1/admin/ingest/stocks \\")
	fmt.Println("     -H 'Authorization: Bearer " + adminToken + "'")
}

func generateJWTSecret() string {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		fmt.Printf("Error generating random secret: %v\n", err)
		os.Exit(1)
	}
	return base64.StdEncoding.EncodeToString(secret)
}

func generateAdminToken(jwtSecret string) string {
	secretBytes, err := base64.StdEncoding.DecodeString(jwtSecret)
	if err != nil {
		fmt.Printf("Error decoding JWT_SECRET: %v\n", err)
		os.Exit(1)
	}

	claims := jwt.MapClaims{
		"sub":  "admin",
		"role": "admin",
		"iat":  time.Now().Unix(),
		"exp":  time.Now().AddDate(0, 1, 0).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretBytes)
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		os.Exit(1)
	}

	return tokenString
}

func updateEnvFile(jwtSecret, adminToken string) {
	envPath := ".env"

	content := ""
	if _, err := os.Stat(envPath); err == nil {
		bytes, err := os.ReadFile(envPath)
		if err != nil {
			fmt.Printf("Error reading .env file: %v\n", err)
			os.Exit(1)
		}
		content = string(bytes)
	}

	content = updateEnvVariable(content, "JWT_SECRET", jwtSecret)

	content = updateEnvVariable(content, "ADMIN_TOKEN", adminToken)

	err := os.WriteFile(envPath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error writing .env file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated %s with JWT_SECRET and ADMIN_TOKEN\n", envPath)
}

func updateEnvVariable(content, key, value string) string {
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			lines[i] = key + "=" + value
			return strings.Join(lines, "\n")
		}
	}

	if len(lines) > 0 && lines[len(lines)-1] != "" {
		lines = append(lines, "")
	}
	lines = append(lines, key+"="+value)

	return strings.Join(lines, "\n")
}
