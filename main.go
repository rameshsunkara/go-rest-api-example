package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bogdanutanu/go-rest-api-example/internal/config"
	"github.com/bogdanutanu/go-rest-api-example/internal/server"
	"github.com/bogdanutanu/go-rest-api-example/internal/utilities"
	"github.com/bogdanutanu/go-rest-api-example/pkg/logger"
	"github.com/bogdanutanu/go-rest-api-example/pkg/mongodb"
	"github.com/rs/zerolog"
)

const (
	serviceName = "ecommerce-orders"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Service %s exited with error: %v (exit code: %d)\n",
			serviceName, err, exitCode(err))
		os.Exit(exitCode(err))
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// setup : read environmental configurations
	svcEnv, envErr := config.Load()
	if envErr != nil {
		return envErr
	}

	// setup : service logger
	var logWriter io.Writer = os.Stdout
	if utilities.IsDevMode(svcEnv.Environment) {
		logWriter = zerolog.ConsoleWriter{Out: os.Stdout}
	}
	lgr := logger.New(svcEnv.LogLevel, logWriter)

	// setup : database connection
	dbConnMgr, dbErr := setupDB(svcEnv)
	if dbErr != nil {
		return dbErr
	}

	lgr.Info().
		Str("name", serviceName).
		Str("environment", svcEnv.Environment).
		Str("started at", time.Now().UTC().Format(time.RFC3339)).
		Msg("Starting the service")

	// Start server - this blocks until shutdown or error
	err := server.Start(ctx, svcEnv, lgr, dbConnMgr)

	// Cleanup after server stops
	cleanup(lgr, dbConnMgr)

	// Don't treat context cancellation as an error
	if errors.Is(err, context.Canceled) {
		lgr.Info().Msg("Graceful shutdown completed")
		return nil
	}

	return err
}

func setupDB(svcEnv *config.ServiceEnvConfig) (*mongodb.ConnectionManager, error) {
	dbCredentials, err := mongodb.CredentialFromSideCar(svcEnv.DBCredentialsSideCar)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch DB credentials : %w", err)
	}

	dbConnMgr, dbErr := mongodb.NewMongoManager(
		svcEnv.DBHosts,
		svcEnv.DBName,
		dbCredentials,
		mongodb.WithQueryLogging(svcEnv.DBLogQueries),
		// mongodb.WithReplicaSet(svcEnv.ReplicaSet) added to demonstrate functional options
	)
	if dbErr != nil {
		return nil, fmt.Errorf("unable to initialize DB connection: %w", dbErr)
	}
	return dbConnMgr, nil
}

func cleanup(lgr logger.Logger, dbConnMgr *mongodb.ConnectionManager) {
	lgr.Info().Msg("Cleaning up resources")
	if err := dbConnMgr.Disconnect(); err != nil {
		lgr.Error().Err(err).Msg("failed to close DB connection, potential connection leak")
		return
	}
}

func exitCode(err error) int {
	if err == nil || errors.Is(err, context.Canceled) {
		return 0
	}
	return 1
}
