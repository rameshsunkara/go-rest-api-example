package models

type ServiceEnv struct {
	Name              string // name of environment where this service is running
	Port              string // port on which this service runs, defaults to DefaultPort
	DBName            string // name of the database
	PrintQueries      bool   // should we print the DB queries that are triggered through this service, defaults to false
	MongoVaultSideCar string // path to find the mongo sidecar file
	DisableAuth       bool   // disables authentication, added to make local development/testing easy
	LogLevel          string // logger level for the service
}
