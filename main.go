package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rameshsunkara/go-rest-api-example/internal/config"
	"github.com/rameshsunkara/go-rest-api-example/internal/server"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
	"github.com/rameshsunkara/go-rest-api-example/pkg/logger"
	"github.com/rameshsunkara/go-rest-api-example/pkg/mongodb"
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
	errCh := make(chan error, 1)

	// setup : read environmental configurations
	svcEnv, envErr := config.Load()
	if envErr != nil {
		return envErr
	}

	// setup : service logger
	var logWriter io.Writer = os.Stdout
	if util.IsDevMode(svcEnv.Environment) {
		logWriter = zerolog.ConsoleWriter{Out: os.Stdout}
	}
	lgr := logger.New(svcEnv.LogLevel, logWriter)

	// setup : database connection
	dbConnMgr, dbErr := setupDB(lgr, svcEnv)
	if dbErr != nil {
		return dbErr
	}

	// setup : create router
	router, routerErr := server.WebRouter(svcEnv, lgr, dbConnMgr)
	if routerErr != nil {
		return routerErr
	}

	// Log registered routes
	lgr.Info().Msg("Registered routes")
	for _, item := range router.Routes() {
		lgr.Info().Str("method", item.Method).Str("path", item.Path).Send()
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + svcEnv.Port,
		Handler: router,
	}

	// Start http server in its own go routine
	go func() {
		lgr.Info().
			Str("name", serviceName).
			Str("environment", svcEnv.Environment).
			Str("started at", time.Now().UTC().Format(time.RFC3339)).
			Msg("Starting the service")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Wait until termination or a critical error
	select {
	case <-ctx.Done():
		lgr.Info().Msg("graceful shutdown signal received")

		// Shutdown server gracefully with 5 second timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			lgr.Error().Err(err).Msg("Server forced to shutdown")
		} else {
			lgr.Info().Msg("Server shutdown gracefully")
		}

		cleanup(lgr, dbConnMgr)
		return nil
	case err := <-errCh:
		lgr.Error().Err(err).Msg("Server error occurred")
		cleanup(lgr, dbConnMgr)
		return err
	}
}

func setupDB(lgr logger.Logger, svcEnv *config.ServiceEnvConfig) (*mongodb.ConnectionManager, error) {
	dbCredentials, err := mongodb.MongoDBCredentialFromSideCar(svcEnv.DBCredentialsSideCar)
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
