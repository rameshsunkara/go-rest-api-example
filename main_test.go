package main

import (
	"context"
	"errors"
	"testing"

	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetEnv(t *testing.T) {
	t.Setenv("environment", "")
	t.Setenv("port", "")
	t.Setenv("dbHosts", "")
	t.Setenv("dbName", "")
	t.Setenv("DBCredentialsSideCar", "")
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
		t.Setenv("dbHosts", "localhost")
		t.Setenv("dbName", "testDB")
		t.Setenv("DBCredentialsSideCar", "/path/to/mongo/sidecar")
		t.Setenv("logLevel", "debug") // Call getEnvConfig and verify returned configurations
		expectedConfig := &models.ServiceEnvConfig{
			Environment:          "test",
			Port:                 "8080",
			LogLevel:             "debug",
			DBCredentialsSideCar: "/path/to/mongo/sidecar",
			DBHosts:              "localhost",
			DBName:               "testDB",
			DBLogQueries:         false, // default value
			DisableAuth:          false, // default value
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
		expectedConfig := &models.ServiceEnvConfig{
			Environment:          "test",
			Port:                 defaultPort,
			LogLevel:             "info",
			DBCredentialsSideCar: "/path/to/mongo/sidecar",
			DBHosts:              "localhost",
			DBName:               "testDB",
			DBLogQueries:         false, // default value
			DisableAuth:          false, // default value
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

func Test_exitCode(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 0, exitCode(nil))
	assert.Equal(t, 1, exitCode(errors.New("fail")))
	assert.Equal(t, 0, exitCode(context.Canceled))
}

func Test_setupDB_fail(t *testing.T) {
	t.Parallel()
	lgr := logger.Setup("info", "test")
	svcEnv := &models.ServiceEnvConfig{
		DBHosts:              "localhost",
		DBName:               "db",
		DBCredentialsSideCar: "/notfound",
	}
	_, err := setupDB(lgr, svcEnv)
	require.Error(t, err)
}
