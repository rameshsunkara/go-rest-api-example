package main

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/bogdanutanu/go-rest-api-example/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExitCode(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 0, exitCode(nil))
	assert.Equal(t, 1, exitCode(errors.New("fail")))
	assert.Equal(t, 0, exitCode(context.Canceled))
}

func TestSetupDBFail(t *testing.T) {
	t.Parallel()
	svcEnv := &config.ServiceEnvConfig{
		DBHosts:              "localhost",
		DBName:               "db",
		DBCredentialsSideCar: "/notfound",
	}
	_, err := setupDB(svcEnv)
	require.Error(t, err)
}

func TestServiceName(t *testing.T) {
	t.Parallel()

	// Test the service name constant
	assert.Equal(t, "ecommerce-orders", serviceName)
}

func TestSetupDBSuccess(t *testing.T) {
	t.Parallel()

	// Create temporary credentials file
	tempFile, err := os.CreateTemp("", "test-credentials-*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write valid credentials
	_, err = tempFile.WriteString(`{"username": "test", "password": "test"}`)
	require.NoError(t, err)
	tempFile.Close()

	svcEnv := &config.ServiceEnvConfig{
		DBHosts:              "mongodb://invalidhost:27017", // Invalid host to avoid actual connection
		DBName:               "testdb",
		DBCredentialsSideCar: tempFile.Name(),
		DBLogQueries:         false,
	}

	// This should fail due to invalid host, but test the credential loading part works
	_, err = setupDB(svcEnv)
	require.Error(t, err)

	// Ensure error is about connection, not credential loading
	assert.Contains(t, err.Error(), "unable to initialize DB connection")
}
