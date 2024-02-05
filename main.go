package main

import (
	"os"
	"time"

	"github.com/rameshsunkara/deferrun"
	"github.com/rameshsunkara/go-rest-api-example/internal/server"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/types"
	"github.com/rs/zerolog/log"
)

const (
	ServiceName = "ecommerce-orders"
	DBName      = "ecommerce"
)

// Passed while building from  make file
var version string

func main() {
	upTime := time.Now()
	t := deferrun.NewSignalHandler()

	env := os.Getenv("environment")
	if env == "" {
		env = "dev"
	}

	// Metadata of the service
	serviceInfo := &types.ServiceInfo{
		Name:        ServiceName,
		UpTime:      upTime,
		Environment: env,
		Version:     version,
	}

	// Setup : Log
	setupLog(env)

	log.Info().Object("Service", serviceInfo).Msg("starting")

	// Setup : DB
	dbManager, dErr := db.NewMongoManager(DBName, "")
	if dErr != nil {
		log.Fatal().Err(dErr).Msg("unable to initialize DB connection")
	}
	t.OnSignal(func() {
		err := dbManager.Disconnect()
		if err != nil {
			log.Err(err).Msg("unable to disconnect from DB")
			return
		}
	})

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
