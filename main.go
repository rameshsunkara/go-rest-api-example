package main

import (
	"os"
	"strconv"
	"time"

	"github.com/rameshsunkara/deferrun"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/server"
	"github.com/rameshsunkara/go-rest-api-example/internal/types"
)

const (
	serviceName = "ecommerce-orders"
	defaultPort = "8080"
)

// Passed while building from  make file.
var version string

func main() {
	upTime := time.Now().UTC().Format(time.RFC3339)
	sigHandler := deferrun.NewSignalHandler()

	// read all environmental configurations, panics if something critical is missing
	svcEnv := MustEnvConfig()

	// service details that can be shared internally for debugging
	svcInfo := types.ServiceInfo{
		Name:        serviceName,
		UpTime:      upTime,
		Environment: svcEnv.Name,
		Version:     version,
	}

	// setup : logger
	lgr := logger.New(svcEnv.Name)

	lgr.ZLog.Info().Object("serviceDetails", svcInfo).Msg("starting")

	// setup : database connection
	dbCredentials, err := db.MongoDBCredentialFromSideCar(svcEnv.MongoVaultSideCar)
	if err != nil {
		lgr.ZLog.Fatal().Err(err).Msg("failed to fetch DB credentials")
	}
	opts := &db.ConnectionOpts{
		Database:     svcEnv.DBName,
		PrintQueries: svcEnv.PrintQueries,
	}
	connMgr, err := db.NewMongoManager(dbCredentials, opts, lgr.ZLog)
	if err != nil {
		lgr.ZLog.Fatal().Err(err).Msg("unable to initialize DB connection")
	}
	sigHandler.OnSignal(func() {
		dErr := connMgr.Disconnect()
		if dErr != nil {
			lgr.ZLog.Err(dErr).Msg("unable to disconnect from DB, potential connection leak")
			return
		}
	})

	// setup : start service - blocking call
	server.StartService(svcInfo, svcEnv, connMgr, lgr)

	lgr.ZLog.Fatal().Object("serviceDetails", svcInfo).Msg("server exited")
}

func MustEnvConfig() types.ServiceEnv {
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

	printDBQueries, err := strconv.ParseBool(os.Getenv("printDBQueries"))
	if err != nil {
		printDBQueries = false
	}

	mongoSideCar := os.Getenv("MongoVaultSideCar")
	if mongoSideCar == "" {
		panic("mongo sidecar file path should be defined in env configuration")
	}

	disableAuth, authEnvErr := strconv.ParseBool(os.Getenv("disableAuth"))
	if authEnvErr != nil {
		// do not disable authentication by default, added this flexibility just for local development purpose
		disableAuth = false
	}

	envConfigurations := types.ServiceEnv{
		Name:              envName,
		Port:              port,
		PrintQueries:      printDBQueries,
		MongoVaultSideCar: mongoSideCar,
		DisableAuth:       disableAuth,
		DBName:            dbName,
	}

	return envConfigurations
}
