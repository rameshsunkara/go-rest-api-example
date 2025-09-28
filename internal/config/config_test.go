package config_test

import (
	"testing"

	"github.com/rameshsunkara/go-rest-api-example/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadWithValidConfig(t *testing.T) {
	// Set required environment variables
	t.Setenv("dbHosts", "localhost:27017")
	t.Setenv("DBCredentialsSideCar", "/path/to/credentials")

	cfg, err := config.Load()

	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "localhost:27017", cfg.DBHosts)
	assert.Equal(t, "/path/to/credentials", cfg.DBCredentialsSideCar)
}

func TestLoadMissingRequiredConfig(t *testing.T) {
	// Clear environment variables
	t.Setenv("dbHosts", "")
	t.Setenv("DBCredentialsSideCar", "")

	_, err := config.Load()

	assert.Error(t, err)
}

func TestLoadWithOptionalDefaults(t *testing.T) {
	// Set only required variables
	t.Setenv("dbHosts", "localhost:27017")
	t.Setenv("DBCredentialsSideCar", "/path/to/credentials")

	// Clear optional variables to test defaults
	t.Setenv("environment", "")
	t.Setenv("port", "")
	t.Setenv("logLevel", "")

	cfg, err := config.Load()

	require.NoError(t, err)
	assert.NotNil(t, cfg)
	// Test that defaults are applied (actual values will depend on the implementation)
	assert.NotEmpty(t, cfg.Environment)
	assert.NotEmpty(t, cfg.Port)
	assert.NotEmpty(t, cfg.LogLevel)
}

func TestConstants(t *testing.T) {
	// Test that constants are properly defined
	assert.Equal(t, "local", config.DefEnvironment)
	assert.Equal(t, "8080", config.DefaultPort)
	assert.Equal(t, "info", config.DefaultLogLevel)
	assert.Equal(t, "ecommerce", config.DefDatabase)
	assert.False(t, config.DefDBQueryLogging)
}

func TestServiceEnvConfigStruct(t *testing.T) {
	// Test that we can create the struct
	cfg := &config.ServiceEnvConfig{
		Environment:          "test",
		Port:                 "8080",
		LogLevel:             "debug",
		DBHosts:              "localhost:27017",
		DBName:               "testdb",
		DBCredentialsSideCar: "/test/path",
		DBLogQueries:         true,
		DisableAuth:          false,
	}

	assert.Equal(t, "test", cfg.Environment)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, "localhost:27017", cfg.DBHosts)
	assert.Equal(t, "testdb", cfg.DBName)
	assert.Equal(t, "/test/path", cfg.DBCredentialsSideCar)
	assert.True(t, cfg.DBLogQueries)
	assert.False(t, cfg.DisableAuth)
}
