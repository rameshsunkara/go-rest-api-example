package main

import (
	"os"
	"strconv"
	"time"

	"github.com/rameshsunkara/deferrun"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/rameshsunkara/go-rest-api-example/internal/server"
)

const (
	serviceName = "ecommerce-orders"
	defaultPort = "8080"
)

// Passed while building from  make file.
var version string

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	upTime := time.Now().UTC().Format(time.RFC3339)
	sigHandler := deferrun.NewSignalHandler()

	// setup : read environmental configurations
	svcEnv := MustEnvConfig()

	// setup : service logger
	lgr := logger.Setup(svcEnv)

	// setup : database connection
	dbCredentials, err := db.MongoDBCredentialFromSideCar(svcEnv.MongoVaultSideCar)
	if err != nil {
		lgr.Fatal().Err(err).Msg("failed to fetch DB credentials")
		return err
	}
	connOpts := &db.ConnectionOpts{
		Database:     svcEnv.DBName,
		PrintQueries: svcEnv.PrintQueries,
	}
	dbConnMgr, err := db.NewMongoManager(dbCredentials, connOpts, lgr)
	if err != nil {
		lgr.Fatal().Err(err).Msg("unable to initialize DB connection")
		return err
	}
	sigHandler.OnSignal(func() {
		dErr := dbConnMgr.Disconnect()
		if dErr != nil {
			lgr.Error().Err(dErr).Msg("unable to disconnect from DB, potential connection leak")
			return
		}
	})

	lgr.Info().
		Str("name", serviceName).
		Str("environment", svcEnv.Name).
		Str("started", upTime).
		Str("version", version).
		Msg("service details, starting the service")

	// setup : start service
	server.StartService(svcEnv, dbConnMgr, lgr)

	lgr.Fatal().Msg("service stopped")
	return nil
}

// MustEnvConfig reads all the environmental configurations and panics if something critical is missing.
func MustEnvConfig() models.ServiceEnv {
	envName := os.Getenv("environment")
	if envName == "" {
		envName = "local"
	}

	port := os.Getenv("port")
	if port == "" {
		port = defaultPort
	}

	dbName := os.Getenv("dbName")
	if dbName == "" {
		panic("dbName should be defined in env configuration")
	}

	// printDBQueries is optional, default is false, when set to true, it will print all the queries to the console.
	printDBQueries, err := strconv.ParseBool(os.Getenv("printDBQueries"))
	if err != nil {
		printDBQueries = false
	}

	mongoSideCar := os.Getenv("MongoVaultSideCar")
	if mongoSideCar == "" {
		panic("mongo sidecar file path should be defined in env configuration")
	}

	// disableAuth is optional, default is false, when set to true, it will disable authentication.
	// Added for development purpose, do not use in production.
	disableAuth, authEnvErr := strconv.ParseBool(os.Getenv("disableAuth"))
	if authEnvErr != nil {
		// do not disable authentication by default, added this flexibility just for local development purpose
		disableAuth = false
	}

	logLevel := os.Getenv("logLevel")
	if logLevel == "" {
		logLevel = "info"
	}

	envConfigurations := models.ServiceEnv{
		Name:              envName,
		Port:              port,
		PrintQueries:      printDBQueries,
		MongoVaultSideCar: mongoSideCar,
		DisableAuth:       disableAuth,
		DBName:            dbName,
		LogLevel:          logLevel,
	}

	return envConfigurations
}
