package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type TestConfig struct {
	DatabaseURLTest string
}

func LoadTestConfig() *TestConfig {
	if err := godotenv.Load("../../.env.test"); err != nil {
		logrus.Warn("No .env.test file found, using environment variables")
	}

	TestConfig := &TestConfig{
		DatabaseURLTest: os.Getenv("DATABASE_URL_TEST"),
	}

	return TestConfig
}

func (tc *TestConfig) HasTestDatabase() bool {
	return tc.DatabaseURLTest != ""
}

func (tc *TestConfig) GetTestDatabaseURL() string {
	return tc.DatabaseURLTest
}
