package logger_test

import (
	"bytes"
	"testing"

	"github.com/rameshsunkara/go-rest-api-example/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	log := logger.New("info", buf)

	assert.NotNil(t, log)
}

func TestDebugLogging(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	log := logger.New("debug", buf)

	log.Debug().Msg("debug message")
	assert.Contains(t, buf.String(), "debug message")
}

func TestInfoLogging(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	log := logger.New("info", buf)

	log.Info().Msg("info message")
	assert.Contains(t, buf.String(), "info message")
}

func TestErrorLogging(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	log := logger.New("error", buf)

	log.Error().Msg("error message")
	assert.Contains(t, buf.String(), "error message")
}

func TestStringField(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	log := logger.New("info", buf)

	log.Info().Str("key", "value").Msg("test")
	output := buf.String()
	assert.Contains(t, output, "key")
	assert.Contains(t, output, "value")
}

func TestIntField(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	log := logger.New("info", buf)

	log.Info().Int("number", 42).Msg("test")
	output := buf.String()
	assert.Contains(t, output, "42")
}

func TestDefaultRequestIDKey(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "X-Request-ID", logger.DefaultRequestIDKey)
}
