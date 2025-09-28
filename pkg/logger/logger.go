package logger

import (
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	DefaultRequestIDKey = "X-Request-ID"
)

// Logger defines the logging interface with method chaining support.
type Logger interface {
	Debug() Event
	Info() Event
	Error() Event
	Fatal() Event
	WithReqID(ctx *gin.Context) (Logger, string)
	WithReqIDCustom(ctx *gin.Context, identifier string) (Logger, string)
}

// Event defines the interface for log event building with method chaining.
type Event interface {
	Str(key, val string) Event
	Int(key string, val int) Event
	Interface(key string, val interface{}) Event
	Dur(key string, val time.Duration) Event
	Err(err error) Event
	Msg(msg string)
	Send()
}

// AppLogger is a zerolog-based implementation of Logger.
type AppLogger struct {
	zLogger zerolog.Logger
}

func New(logLevel string, writer io.Writer) Logger {
	lvl := parseZerologLevel(logLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano
	return &AppLogger{
		zLogger: zerolog.New(writer).With().Caller().Timestamp().Logger().Level(lvl),
	}
}

// WithReqID returns a logger with request ID using DefaultRequestIDKey.
func (l *AppLogger) WithReqID(ctx *gin.Context) (Logger, string) {
	return l.WithReqIDCustom(ctx, DefaultRequestIDKey)
}

// WithReqIDCustom returns a logger with request ID using a custom identifier key.
func (l *AppLogger) WithReqIDCustom(ctx *gin.Context, identifier string) (Logger, string) {
	idKey := identifier

	type contextKey string
	if rID := ctx.Request.Context().Value(contextKey(idKey)); rID != nil {
		if reqID, ok := rID.(string); ok {
			return &AppLogger{zLogger: l.zLogger.With().Str(idKey, reqID).Logger()}, reqID
		}
		return l, ""
	}
	return l, ""
}

// zerologEvent wraps zerolog.Event to implement the Event interface.
type zerologEvent struct {
	event *zerolog.Event
}

// Implement Event interface
func (e *zerologEvent) Str(key, val string) Event {
	return &zerologEvent{event: e.event.Str(key, val)}
}

func (e *zerologEvent) Int(key string, val int) Event {
	return &zerologEvent{event: e.event.Int(key, val)}
}

func (e *zerologEvent) Interface(key string, val interface{}) Event {
	return &zerologEvent{event: e.event.Interface(key, val)}
}

func (e *zerologEvent) Dur(key string, val time.Duration) Event {
	return &zerologEvent{event: e.event.Dur(key, val)}
}

func (e *zerologEvent) Err(err error) Event {
	return &zerologEvent{event: e.event.Err(err)}
}

func (e *zerologEvent) Msg(msg string) {
	e.event.Msg(msg)
}

func (e *zerologEvent) Send() {
	e.event.Send()
}

// Fatal logs a message with fatal level and exits the program.
func (l *AppLogger) Fatal() Event {
	return &zerologEvent{event: l.zLogger.Fatal()}
}

// Error logs a message with error level.
func (l *AppLogger) Error() Event {
	return &zerologEvent{event: l.zLogger.Error()}
}

// Info logs a message with info level.
func (l *AppLogger) Info() Event {
	return &zerologEvent{event: l.zLogger.Info()}
}

// Debug logs a message with debug level.
func (l *AppLogger) Debug() Event {
	return &zerologEvent{event: l.zLogger.Debug()}
}

// parseZerologLevel parses a string log level to zerolog.Level.
func parseZerologLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
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
