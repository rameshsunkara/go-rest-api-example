package logger

import (
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var once sync.Once
var logger zerolog.Logger

func SetupZeroLogger(env string) *zerolog.Logger {
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
	return &logger
}

func ZeroLogger() *zerolog.Logger {
	return &logger
}

const RequestIdentifier = "X-Request-ID"

func ReqLogger(c *gin.Context) (zerolog.Logger, string) {
	reqContext := c.Request.Context()
	if rId := reqContext.Value(RequestIdentifier); rId != nil {
		reqId := rId.(string)
		return logger.With().Str(RequestIdentifier, reqId).Logger(), reqId
	}
	return logger, ""
}
