package config

import (
	"errors"
	"os"
	"strconv"
)

// ServiceEnvConfig holds all environmental configurations for the service.
type ServiceEnvConfig struct {
	Environment string // environment where this service is running (dev, staging, prod, etc.)
	Port        string // port on which this service runs, defaults to DefaultPort
	LogLevel    string // logger level for the service

	// DB related configurations
	DBCredentialsSideCar string // path to find the database credentials sidecar file
	DBHosts              string // comma separated list of DB hosts
	DBName               string // name of the database
	DBPort               int    // port on which the DB is listening, defaults to 27017
	DBLogQueries         bool   // print the DB queries that are triggered through this service, defaults to false

	DisableAuth   bool // disables API authentication, added to make local development/testing easy
	EnableTracing bool // enables flight recorder for slow request tracing, defaults to false
}

const (
	DefaultPort       = "8080"
	DefaultLogLevel   = "info"
	DefDatabase       = "ecommerce"
	DefEnvironment    = "local"
	DefDBQueryLogging = false
	DefEnableTracing  = false
)

// Load reads all environmental configurations and returns a ServiceEnvConfig.
func Load() (*ServiceEnvConfig, error) {
	dbCredentialsSideCar := os.Getenv("DBCredentialsSideCar")
	if dbCredentialsSideCar == "" {
		return nil, errors.New("database credentials sidecar file path is missing in env")
	}

	envName := os.Getenv("environment")
	if envName == "" {
		envName = DefEnvironment
	}

	port := os.Getenv("port")
	if port == "" {
		port = DefaultPort
	}

	dbHosts := os.Getenv("dbHosts")
	if dbHosts == "" {
		return nil, errors.New("dbHosts is missing in env")
	}
	dbName := os.Getenv("dbName")
	if dbName == "" {
		dbName = DefDatabase
	}
	printDBQueries, err := strconv.ParseBool(os.Getenv("printDBQueries"))
	if err != nil {
		printDBQueries = DefDBQueryLogging
	}

	disableAuth, authEnvErr := strconv.ParseBool(os.Getenv("disableAuth"))
	if authEnvErr != nil {
		// do not disable authentication by default, added this flexibility just for local development purpose
		disableAuth = false
	}

	enableTracing, tracingEnvErr := strconv.ParseBool(os.Getenv("enableTracing"))
	if tracingEnvErr != nil {
		// tracing is disabled by default to avoid overhead in production
		enableTracing = DefEnableTracing
	}

	logLevel := os.Getenv("logLevel")
	if logLevel == "" {
		logLevel = DefaultLogLevel
	}

	envConfigurations := &ServiceEnvConfig{
		Environment:          envName,
		Port:                 port,
		DBHosts:              dbHosts,
		DBName:               dbName,
		DBCredentialsSideCar: dbCredentialsSideCar,
		DisableAuth:          disableAuth,
		EnableTracing:        enableTracing,
		LogLevel:             logLevel,
		DBLogQueries:         printDBQueries,
	}

	return envConfigurations, nil
}
