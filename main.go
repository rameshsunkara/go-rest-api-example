package main

import (
	"os"
	"time"

	"github.com/rameshsunkara/go-rest-api-example/internal/server"
	"github.com/rameshsunkara/go-rest-api-example/pkg/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	_ "github.com/rameshsunkara/go-rest-api-example/docs"
	"github.com/rameshsunkara/go-rest-api-example/internal/config"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/rs/zerolog/log"
)

const (
	ServiceName = "ecommerce-orders"
	DBName      = "ecommerce"
)

// Passed while building from  make file
var version string

// @title           GO Rest Example API Service (Purchase Order Tracker)
// @version         1.0
// @description     A sample service to demonstrate how to develop REST API in golang

// @contact.name    Ramesh Sunkara
// @contact.url
// @contact.email

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	upTime := time.Now()

	env := os.Getenv("environment")
	if env == "" {
		env = "dev"
	}

	// Metadata of the service
	serviceInfo := &models.ServiceInfo{
		Name:        ServiceName,
		UpTime:      upTime,
		Environment: env,
		Version:     version,
	}

	// Setup : Log
	setupLog(env)

	log.Info().Object("Service", serviceInfo).Msg("starting")

	// Load Configuration
	c, cErr := config.LoadConfig(env)
	if cErr != nil {
		log.Fatal().Err(cErr).Msg("unable to read configuration")
	}

	// Setup : DB
	dbManager, dErr := db.Init(DBName, c.GetString("db.dsn"))
	if dErr != nil {
		log.Fatal().Err(dErr).Msg("unable to initialize DB connection")
	}

	// Setup : Server
	server.Init(serviceInfo, dbManager)

	log.Fatal().Str("ServiceName", ServiceName).Msg("Server Exited")
}

func setupLog(env string) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	lvl := zerolog.InfoLevel
	logDest := os.Stdout
	logger := zerolog.New(logDest).With().Caller().Timestamp().Logger()
	if util.IsDevMode(env) {
		lvl = zerolog.TraceLevel
		logger = zerolog.New(zerolog.ConsoleWriter{Out: logDest}).With().Caller().Timestamp().Logger()
	}
	zerolog.SetGlobalLevel(lvl)
	log.Logger = logger
}
