package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSidecarPath = "/path/to/mongo/sidecar"
	testDBHosts     = "localhost:27017"
)

// resetEnv clears all environment variables used by config.Load().
func resetEnv(t *testing.T) {
	t.Helper()
	envVars := []string{
		"environment", "port", "dbHosts", "dbName",
		"DBCredentialsSideCar", "logLevel", "printDBQueries", "disableAuth",
	}
	for _, envVar := range envVars {
		t.Setenv(envVar, "")
	}
}

// setMinimalRequiredEnv sets only the required environment variables.
func setMinimalRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("dbHosts", testDBHosts)
	t.Setenv("DBCredentialsSideCar", testSidecarPath)
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		setupEnv    func(*testing.T)
		expected    *ServiceEnvConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "all values explicitly set",
			setupEnv: func(t *testing.T) {
				resetEnv(t)
				t.Setenv("environment", "production")
				t.Setenv("port", "9090")
				t.Setenv("dbHosts", "db1:27017,db2:27017")
				t.Setenv("dbName", "prodDB")
				t.Setenv("DBCredentialsSideCar", "/prod/credentials")
				t.Setenv("logLevel", "error")
				t.Setenv("printDBQueries", "true")
				t.Setenv("disableAuth", "true")
			},
			expected: &ServiceEnvConfig{
				Environment:          "production",
				Port:                 "9090",
				LogLevel:             "error",
				DBCredentialsSideCar: "/prod/credentials",
				DBHosts:              "db1:27017,db2:27017",
				DBName:               "prodDB",
				DBLogQueries:         true,
				DisableAuth:          true,
			},
		},
		{
			name: "minimal required values with defaults",
			setupEnv: func(t *testing.T) {
				resetEnv(t)
				setMinimalRequiredEnv(t)
			},
			expected: &ServiceEnvConfig{
				Environment:          DefEnvironment,
				Port:                 DefaultPort,
				LogLevel:             DefaultLogLevel,
				DBCredentialsSideCar: testSidecarPath,
				DBHosts:              testDBHosts,
				DBName:               DefDatabase,
				DBLogQueries:         DefDBQueryLogging,
				DisableAuth:          false,
			},
		},
		{
			name: "invalid boolean for printDBQueries defaults to false",
			setupEnv: func(t *testing.T) {
				resetEnv(t)
				setMinimalRequiredEnv(t)
				t.Setenv("printDBQueries", "invalid-boolean")
			},
			expected: &ServiceEnvConfig{
				Environment:          DefEnvironment,
				Port:                 DefaultPort,
				LogLevel:             DefaultLogLevel,
				DBCredentialsSideCar: testSidecarPath,
				DBHosts:              testDBHosts,
				DBName:               DefDatabase,
				DBLogQueries:         DefDBQueryLogging, // defaults to false
				DisableAuth:          false,
			},
		},
		{
			name: "invalid boolean for disableAuth defaults to false",
			setupEnv: func(t *testing.T) {
				resetEnv(t)
				setMinimalRequiredEnv(t)
				t.Setenv("disableAuth", "not-a-boolean")
			},
			expected: &ServiceEnvConfig{
				Environment:          DefEnvironment,
				Port:                 DefaultPort,
				LogLevel:             DefaultLogLevel,
				DBCredentialsSideCar: testSidecarPath,
				DBHosts:              testDBHosts,
				DBName:               DefDatabase,
				DBLogQueries:         DefDBQueryLogging,
				DisableAuth:          false, // defaults to false
			},
		},
		{
			name: "missing DBCredentialsSideCar",
			setupEnv: func(t *testing.T) {
				resetEnv(t)
				t.Setenv("dbHosts", testDBHosts)
				// DBCredentialsSideCar intentionally not set
			},
			expectError: true,
			errorMsg:    "database credentials sidecar file path is missing in env",
		},
		{
			name: "missing dbHosts",
			setupEnv: func(t *testing.T) {
				resetEnv(t)
				t.Setenv("DBCredentialsSideCar", testSidecarPath)
				// dbHosts intentionally not set
			},
			expectError: true,
			errorMsg:    "dbHosts is missing in env",
		},
		{
			name: "empty DBCredentialsSideCar",
			setupEnv: func(t *testing.T) {
				resetEnv(t)
				t.Setenv("dbHosts", testDBHosts)
				t.Setenv("DBCredentialsSideCar", "")
			},
			expectError: true,
			errorMsg:    "database credentials sidecar file path is missing in env",
		},
		{
			name: "empty dbHosts",
			setupEnv: func(t *testing.T) {
				resetEnv(t)
				t.Setenv("DBCredentialsSideCar", testSidecarPath)
				t.Setenv("dbHosts", "")
			},
			expectError: true,
			errorMsg:    "dbHosts is missing in env",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv(t)

			actualConfig, err := Load()

			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Nil(t, actualConfig)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, actualConfig)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	t.Parallel()

	// Verify constants are sensible values
	assert.Equal(t, "8080", DefaultPort)
	assert.Equal(t, "info", DefaultLogLevel)
	assert.Equal(t, "ecommerce", DefDatabase)
	assert.Equal(t, "local", DefEnvironment)
	assert.False(t, DefDBQueryLogging)
}

func TestServiceEnvConfigFieldTypes(t *testing.T) {
	t.Parallel()

	// Test that ServiceEnvConfig has the expected field types
	config := &ServiceEnvConfig{}

	// String fields
	assert.IsType(t, "", config.Environment)
	assert.IsType(t, "", config.Port)
	assert.IsType(t, "", config.LogLevel)
	assert.IsType(t, "", config.DBCredentialsSideCar)
	assert.IsType(t, "", config.DBHosts)
	assert.IsType(t, "", config.DBName)

	// Int field
	assert.IsType(t, 0, config.DBPort)

	// Bool fields
	assert.IsType(t, false, config.DBLogQueries)
	assert.IsType(t, false, config.DisableAuth)
}
