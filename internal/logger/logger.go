package logger

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	setupOnce sync.Once
	appLogger *AppLogger
)

// AppLogger is a wrapper around zerolog.Logger.
type AppLogger struct {
	zLogger zerolog.Logger
}

func Setup(env models.ServiceEnv) *AppLogger {
	setupOnce.Do(func() {
		appLogger = &AppLogger{}
		lvl := ZerologLevel(env.LogLevel)
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano
		var logDest io.Writer
		logDest = os.Stdout
		if util.IsDevMode(env.Name) {
			logDest = zerolog.ConsoleWriter{Out: logDest}
		}
		appLogger.zLogger = zerolog.New(logDest).With().Caller().Timestamp().Logger().Level(lvl)
	})
	return appLogger
}

// WithReqID returns a logger with request ID.
func (l *AppLogger) WithReqID(ctx *gin.Context) (zerolog.Logger, string) {
	if rID := ctx.Request.Context().Value(util.ContextKey(util.RequestIdentifier)); rID != nil {
		if reqID, ok := rID.(string); ok {
			return l.zLogger.With().Str(util.RequestIdentifier, reqID).Logger(), reqID
		}
		return l.zLogger, ""
	}
	return l.zLogger, ""
}

// Fatal logs a message with fatal level and exits the program.
func (l *AppLogger) Fatal() *zerolog.Event {
	return l.zLogger.Fatal()
}

// Error logs a message with error level.
func (l *AppLogger) Error() *zerolog.Event {
	return l.zLogger.Error()
}

// Info logs a message with info level.
func (l *AppLogger) Info() *zerolog.Event {
	return l.zLogger.Info()
}

// Debug logs a message with debug level.
func (l *AppLogger) Debug() *zerolog.Event {
	return l.zLogger.Debug()
}

func ZerologLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
