package logger_test

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	// Prepare mock service environment
	env := models.ServiceEnv{Name: "dev"}

	// Call Setup function
	lgr := logger.Setup(env)

	// Ensure logger is initialized
	assert.NotNil(t, lgr)
}

func TestWithReqID(t *testing.T) {
	// Prepare mock service environment
	env := models.ServiceEnv{Name: "test"}

	// Call Setup function
	lgr := logger.Setup(env)

	// Prepare a mock gin context
	ginCtx := &gin.Context{}
	ginCtx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	// Call WithReqID with a context without request ID
	_, reqID := lgr.WithReqID(ginCtx)
	assert.Empty(t, reqID)

	// add a request ID to the context
	reqIDValue := "1234567890"
	ctx := context.WithValue(ginCtx.Request.Context(), util.ContextKey(util.RequestIdentifier), reqIDValue)
	ginCtx.Request = ginCtx.Request.WithContext(ctx)

	// Call WithReqID with a context containing a request ID
	_, newReqID := lgr.WithReqID(ginCtx)
	assert.Equal(t, reqIDValue, newReqID)

	// Call WithReqID with a context containing a non string request ID
	ctx = context.WithValue(ginCtx.Request.Context(), util.ContextKey(util.RequestIdentifier), 123)
	ginCtx.Request = ginCtx.Request.WithContext(ctx)
	_, newReqID = lgr.WithReqID(ginCtx)
	assert.Empty(t, newReqID)
}

func TestSetupOnce(t *testing.T) {
	// Prepare mock service environment
	env := models.ServiceEnv{Name: "test"}

	// Use a temporary file for logging
	tempFile, err := os.CreateTemp("", "uTest.log")
	require.NoError(t, err)
	defer func(name string) {
		errRemove := os.Remove(name)
		if err != nil {
			t.Log(errRemove)
		}
	}(tempFile.Name())

	// Call Setup function concurrently multiple times
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			lgr := logger.Setup(env)
			assert.NotNil(t, lgr)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestGetZerologLevel(t *testing.T) {
	tests := []struct {
		name       string
		inputLevel string
		expected   zerolog.Level
	}{
		{"Debug", "debug", zerolog.DebugLevel},
		{"Info", "info", zerolog.InfoLevel},
		{"Error", "error", zerolog.ErrorLevel},
		{"Fatal", "fatal", zerolog.FatalLevel},
		{"Unknown", "unknown", zerolog.InfoLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := logger.ZerologLevel(tt.inputLevel)
			assert.Equal(t, tt.expected, actual, "Unexpected log level")
		})
	}
}
