package logger

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type Logger struct {
	ZLog zerolog.Logger
}

func New(env string) *Logger {
	al := Logger{}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logDest := os.Stdout
	al.ZLog = zerolog.New(logDest).With().Caller().Timestamp().Logger().Level(zerolog.InfoLevel)

	if util.IsDevMode(env) {
		al.ZLog = zerolog.New(zerolog.ConsoleWriter{Out: logDest}).With().
			Caller().Timestamp().Logger().Level(zerolog.TraceLevel)
	}
	return &al
}

func (l *Logger) WithReqID(ctx *gin.Context) (zerolog.Logger, string) {
	reqContext := ctx.Request.Context()
	if rID := reqContext.Value(util.RequestIdentifier); rID != nil {
		reqID := rID.(string)
		return l.ZLog.With().Str(util.RequestIdentifier, reqID).Logger(), reqID
	}
	return l.ZLog, ""
}
