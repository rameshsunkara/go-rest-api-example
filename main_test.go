package main

import (
	"testing"

	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestMustEnvConfig(t *testing.T) {
	t.Run("MissingEnvVariables", func(t *testing.T) {
		// Call MustEnvConfig and expect panic
		assert.Panics(t, func() {
			_ = MustEnvConfig()
		}, "MustEnvConfig did not panic with missing environment variables")
	})

	t.Run("ValidEnvVariables", func(t *testing.T) {
		// Set required environment variables
		t.Setenv("environment", "test")
		t.Setenv("port", "8080")
		t.Setenv("dbName", "testDB")
		t.Setenv("MongoVaultSideCar", "/path/to/mongo/sidecar")
		t.Setenv("logLevel", "debug")

		// Call MustEnvConfig and verify returned configurations
		expectedConfig := models.ServiceEnv{
			Name:              "test",
			Port:              "8080",
			PrintQueries:      false, // default value
			MongoVaultSideCar: "/path/to/mongo/sidecar",
			DisableAuth:       false, // default value
			DBName:            "testDB",
			LogLevel:          "debug",
		}

		actualConfig := MustEnvConfig()
		assert.Equal(t, expectedConfig, actualConfig,
			"MustEnvConfig did not return expected configurations")
	})
}

func TestMustEnvConfig_Defaults(t *testing.T) {
	t.Run("Default Env Values", func(t *testing.T) {
		// Set required environment variables
		t.Setenv("environment", "test")
		t.Setenv("dbName", "testDB")
		t.Setenv("MongoVaultSideCar", "/path/to/mongo/sidecar")

		// Call MustEnvConfig and verify returned configurations
		expectedConfig := models.ServiceEnv{
			Name:              "test",
			Port:              defaultPort,
			PrintQueries:      false, // default value
			MongoVaultSideCar: "/path/to/mongo/sidecar",
			DisableAuth:       false, // default value
			DBName:            "testDB",
			LogLevel:          "info",
		}

		actualConfig := MustEnvConfig()
		assert.Equal(t, expectedConfig, actualConfig,
			"Default values are not set correctly in MustEnvConfig")
	})
}

func TestMustEnvConfig_FailOnSideCar(t *testing.T) {
	t.Run("Fail on MongoSide Car", func(t *testing.T) {
		// Set required environment variables
		t.Setenv("environment", "test")
		t.Setenv("dbName", "testDB")

		assert.Panics(t, func() {
			_ = MustEnvConfig()
		}, "MustEnvConfig did not panic with missing environment variables")
	})
}
