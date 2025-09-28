package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/rameshsunkara/go-rest-api-example/internal/server"
)

const (
	serviceName       = "ecommerce-orders"
	defaultPort       = "8080"
	defaultLogLevel   = "info"
	defDatabase       = "ecommerce"
	defEnvironment    = "local"
	defDBQueryLogging = false
)

// Passed while building from the make file.
var version string

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
	svcEnv, envErr := getEnvConfig()
	if envErr != nil {
		return envErr
	}

	// setup : service logger
	lgr := logger.Setup(svcEnv.LogLevel, svcEnv.Environment)

	// setup : database connection
	dbConnMgr, dbErr := setupDB(lgr, svcEnv)
	if dbErr != nil {
		return dbErr
	}

	// Start http server in its own go routine
	go func() {
		errCh <- server.Start(svcEnv, lgr, dbConnMgr)
	}()

	lgr.Info().
		Str("name", serviceName).
		Str("environment", svcEnv.Environment).
		Str("started at", time.Now().UTC().Format(time.RFC3339)).
		Str("version", version).
		Msg("starting the service")

	// Wait until termination or a critical error
	select {
	case <-ctx.Done():
		lgr.Info().Msg("graceful shutdown signal received")
		err := <-errCh // wait for go routines to exit
		cleanup(lgr, dbConnMgr)
		return err
	case err := <-errCh:
		lgr.Error().Err(err).Msg("something went wrong")
		cleanup(lgr, dbConnMgr)
		return err
	}
}

// getEnvConfig reads all the environmental configurations.
func getEnvConfig() (*models.ServiceEnvConfig, error) {
	dbCredentialsSideCar := os.Getenv("DBCredentialsSideCar")
	if dbCredentialsSideCar == "" {
		return nil, errors.New("database credentials sidecar file path is missing in env")
	}

	envName := os.Getenv("environment")
	if envName == "" {
		envName = defEnvironment
	}

	port := os.Getenv("port")
	if port == "" {
		port = defaultPort
	}

	dbHosts := os.Getenv("dbHosts")
	if dbHosts == "" {
		return nil, errors.New("dbHosts is missing in env")
	}
	dbName := os.Getenv("dbName")
	if dbName == "" {
		dbName = defDatabase
	}
	printDBQueries, err := strconv.ParseBool(os.Getenv("printDBQueries"))
	if err != nil {
		printDBQueries = defDBQueryLogging
	}

	disableAuth, authEnvErr := strconv.ParseBool(os.Getenv("disableAuth"))
	if authEnvErr != nil {
		// do not disable authentication by default, added this flexibility just for local development purpose
		disableAuth = false
	}

	logLevel := os.Getenv("logLevel")
	if logLevel == "" {
		logLevel = defaultLogLevel
	}

	envConfigurations := &models.ServiceEnvConfig{
		Environment:          envName,
		Port:                 port,
		DBHosts:              dbHosts,
		DBName:               dbName,
		DBCredentialsSideCar: dbCredentialsSideCar,
		DisableAuth:          disableAuth,
		LogLevel:             logLevel,
		DBLogQueries:         printDBQueries,
	}

	return envConfigurations, nil
}

func setupDB(lgr *logger.AppLogger, svcEnv *models.ServiceEnvConfig) (*db.ConnectionManager, error) {
	dbCredentials, err := db.MongoDBCredentialFromSideCar(svcEnv.DBCredentialsSideCar)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch DB credentials : %w", err)
	}

	// Parse comma-separated hosts into slice
	dbHosts := strings.Split(strings.ReplaceAll(svcEnv.DBHosts, " ", ""), ",")

	dbConnMgr, dbErr := db.NewMongoManager(
		dbHosts,
		svcEnv.DBName,
		dbCredentials,
		lgr,
		db.WithQueryLogging(svcEnv.DBLogQueries),
		db.WithPort(27022),         // using port 27022 to match docker-compose configuration
		db.WithAuthSource("admin"), // authenticate against admin database for root credentials
	)
	if dbErr != nil {
		return nil, fmt.Errorf("unable to initialize DB connection: %w", dbErr)
	}
	return dbConnMgr, nil
}

func cleanup(lgr *logger.AppLogger, dbConnMgr *db.ConnectionManager) {
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
