package models

type ServiceEnvConfig struct {
	Environment string // environment where this service is running (dev, staging, prod, etc.)
	Port        string // port on which this service runs, defaults to DefaultPort
	LogLevel    string // logger level for the service

	// DB related configurations
	DBCredentialsSideCar string // path to find the database credentials sidecar file
	DBHosts              string // comma separated list of DB hosts
	DBName               string // name of the database
	DBPort               int    // port on which the DB is listening, defaults to 27017
	DBLogQueries         bool   // should we print the DB queries that are triggered through this service, defaults to false

	DisableAuth bool // disables API authentication, added to make local development/testing easy
}
