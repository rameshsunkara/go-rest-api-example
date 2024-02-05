package logger

import (
	"os"
	"sync"

	"github.com/rameshsunkara/go-rest-api-example/internal/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var once sync.Once

func ZeroLogger(env string) zerolog.Logger {
	var logger zerolog.Logger
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		lvl := zerolog.InfoLevel
		logDest := os.Stdout
		logger = zerolog.New(logDest).With().Caller().Timestamp().Logger()

		if util.IsDevMode(env) {
			lvl = zerolog.TraceLevel
			logger = zerolog.New(zerolog.ConsoleWriter{Out: logDest}).With().Caller().Timestamp().Logger()
		}
		zerolog.SetGlobalLevel(lvl)
	})
	return logger
}
