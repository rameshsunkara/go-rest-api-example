package main

import (
	"testing"

	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetEnv(t *testing.T) {
	t.Setenv("environment", "")
	t.Setenv("port", "")
	t.Setenv("dbName", "")
	t.Setenv("MongoVaultSideCar", "")
	t.Setenv("logLevel", "")
	t.Setenv("printDBQueries", "")
}

func TestMustEnvConfig(t *testing.T) {
	t.Run("MissingEnvVariables", func(t *testing.T) {
		resetEnv(t)
		// Call getEnvConfig and expect panic
		_, err := getEnvConfig()
		require.Error(t, err)
	})

	t.Run("ValidEnvVariables", func(t *testing.T) {
		resetEnv(t)
		// Set required environment variables
		t.Setenv("environment", "test")
		t.Setenv("port", "8080")
		t.Setenv("dbName", "testDB")
		t.Setenv("MongoVaultSideCar", "/path/to/mongo/sidecar")
		t.Setenv("logLevel", "debug")

		// Call getEnvConfig and verify returned configurations
		expectedConfig := &models.ServiceEnv{
			Name:              "test",
			Port:              "8080",
			PrintQueries:      false, // default value
			MongoVaultSideCar: "/path/to/mongo/sidecar",
			DisableAuth:       false, // default value
			DBName:            "testDB",
			LogLevel:          "debug",
		}

		actualConfig, err := getEnvConfig()
		assert.Equal(t, expectedConfig, actualConfig,
			"getEnvConfig did not return expected configurations")
		require.NoError(t, err)
	})
}

func TestMustEnvConfig_Defaults(t *testing.T) {
	t.Run("Default Env Values", func(t *testing.T) {
		resetEnv(t)
		// Set required environment variables
		t.Setenv("environment", "test")
		t.Setenv("dbName", "testDB")
		t.Setenv("MongoVaultSideCar", "/path/to/mongo/sidecar")

		// Call getEnvConfig and verify returned configurations
		expectedConfig := &models.ServiceEnv{
			Name:              "test",
			Port:              defaultPort,
			PrintQueries:      false, // default value
			MongoVaultSideCar: "/path/to/mongo/sidecar",
			DisableAuth:       false, // default value
			DBName:            "testDB",
			LogLevel:          "info",
		}

		actualConfig, err := getEnvConfig()
		assert.Equal(t, expectedConfig, actualConfig,
			"Default values are not set correctly in getEnvConfig")
		require.NoError(t, err)
	})
}

func TestMustEnvConfig_FailOnSideCar(t *testing.T) {
	t.Run("Fail on MongoSide Car", func(t *testing.T) {
		resetEnv(t)
		// Set required environment variables
		t.Setenv("environment", "test")
		t.Setenv("dbName", "testDB")

		_, err := getEnvConfig()
		require.Error(t, err)
	})
}
