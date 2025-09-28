package main

import (
	"context"
	"errors"
	"testing"

	"github.com/rameshsunkara/go-rest-api-example/internal/config"
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
