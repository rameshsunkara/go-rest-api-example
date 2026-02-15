package logger_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bogdanutanu/go-rest-api-example/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New("info", buf)

	assert.NotNil(t, log)
}

func TestDebugLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New("debug", buf)

	log.Debug().Msg("debug message")
	assert.Contains(t, buf.String(), "debug message")
}

func TestInfoLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New("info", buf)

	log.Info().Msg("info message")
	assert.Contains(t, buf.String(), "info message")
}

func TestErrorLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New("error", buf)

	log.Error().Msg("error message")
	assert.Contains(t, buf.String(), "error message")
}

func TestStringField(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New("info", buf)

	log.Info().Str("key", "value").Msg("test")
	output := buf.String()
	assert.Contains(t, output, "key")
	assert.Contains(t, output, "value")
}

func TestIntField(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New("info", buf)

	log.Info().Int("number", 42).Msg("test")
	output := buf.String()
	assert.Contains(t, output, "42")
}

func TestDefaultRequestIDKey(t *testing.T) {
	assert.Equal(t, "X-Request-ID", logger.DefaultRequestIDKey)
}

func TestInterfaceField(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New("info", buf)

	testData := map[string]int{"count": 5}
	log.Info().Interface("data", testData).Msg("test")
	output := buf.String()
	assert.Contains(t, output, "data")
	assert.Contains(t, output, "count")
}

func TestDurField(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New("info", buf)

	duration := time.Duration(100) * time.Millisecond
	log.Info().Dur("duration", duration).Msg("test")
	output := buf.String()
	assert.Contains(t, output, "duration")
	assert.Contains(t, output, "100")
}

func TestErrField(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New("info", buf)

	testError := errors.New("test error")
	log.Info().Err(testError).Msg("test")
	output := buf.String()
	assert.Contains(t, output, "test error")
}

func TestSendMethod(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New("info", buf)

	log.Info().Str("key", "value").Send()
	output := buf.String()
	assert.Contains(t, output, "key")
	assert.Contains(t, output, "value")
}

func TestWithReqIDBasic(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New("info", buf)

	// Test that WithReqID doesn't panic and returns the logger
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	loggerWithReqID, reqID := log.WithReqID(c)
	assert.NotNil(t, loggerWithReqID)
	// Without proper context setup, reqID will be empty, but method doesn't panic
	assert.Empty(t, reqID)
}

func TestWithReqIDCustomBasic(t *testing.T) {
	buf := &bytes.Buffer{}
	log := logger.New("info", buf)

	// Test that WithReqIDCustom doesn't panic and returns the logger
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	loggerWithReqID, reqID := log.WithReqIDCustom(c, "Custom-ID")
	assert.NotNil(t, loggerWithReqID)
	// Without proper context setup, reqID will be empty, but method doesn't panic
	assert.Empty(t, reqID)
}
