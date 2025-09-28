package main

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/rameshsunkara/go-rest-api-example/internal/config"
	"github.com/rameshsunkara/go-rest-api-example/pkg/logger"
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
	lgr := logger.New("info", os.Stdout)
	svcEnv := &config.ServiceEnvConfig{
		DBHosts:              "localhost",
		DBName:               "db",
		DBCredentialsSideCar: "/notfound",
	}
	_, err := setupDB(lgr, svcEnv)
	require.Error(t, err)
}
