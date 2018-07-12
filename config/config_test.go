package config

import (
	"os"
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	configVars := map[string]string {
		"username": "postgres",
		"port": "5433",
		"password": "",
		"dbname": "food_test",
	}

	os.Setenv("DEPLOYENV", "ci")
	defer os.Unsetenv("DEPLOYENV")

	username, port, password, dbName := GetConfig()
	assert.Equal(t, configVars["username"], username)
	assert.Equal(t, 5433, port)
	assert.Equal(t, configVars["password"], password)
	assert.Equal(t, configVars["dbname"], dbName)
}
