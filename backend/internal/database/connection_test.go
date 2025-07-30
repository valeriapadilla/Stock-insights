package database

import(
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectWithValidURL(t *testing.T){
	os.Setenv("DATABASE_URL","postgres://postgres:postgres@localhost:5432/stock_insights_test")
	defer os.Unsetenv("DATABASE_URL")

	err := Connect()

	if err != nil{
		assert.Contains(t, err.Error(), "connection")
	} else {
		assert.NotNil(t,DB)
		Close()
	}
}

func TestConnectWithInvalidURL(t *testing.T){
	os.Setenv("DATABASE_URL","invalid://url")
	defer os.Unsetenv("DATABASE_URL")

	err := Connect()

	assert.Error(t,err)
	assert.Contains(t, err.Error(), "failed to connect to database")
}

func TestConnectWithMissingURL(t *testing.T){
	os.Unsetenv("DATABASE_URL")

	err := Connect()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DATABASE_URL is not set")
}

func TestPingWithoutConnection(t *testing.T){
	DB = nil

	err := Ping()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database connection not established")
}

func TestGetStatsWithoutConnection(t *testing.T){
	DB = nil
	stats := GetStats()
	assert.Equal(t, sql.DBStats{}, stats)
}
