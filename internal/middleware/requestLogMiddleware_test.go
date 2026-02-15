package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/bogdanutanu/go-rest-api-example/internal/middleware"
	"github.com/bogdanutanu/go-rest-api-example/pkg/flightrecorder"
	"github.com/bogdanutanu/go-rest-api-example/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRequestLogMiddleware(_ *testing.T) {
	type requestLogMiddlewareTestCase struct {
		Description  string
		InputReqPath string
	}

	var testCases = []requestLogMiddlewareTestCase{
		{
			Description:  "improve assertions-1",
			InputReqPath: "/test/1",
		},
		{
			Description:  "improve assertions-2",
			InputReqPath: "/test/2",
		},
	}

	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(resp)
	lgr := logger.New("info", os.Stdout)
	r.Use(middleware.RequestLogMiddleware(lgr, nil)) // Pass nil for flight recorder in tests

	for _, tc := range testCases {
		r.GET(tc.InputReqPath, func(ctx *gin.Context) {
			ctx.String(200, "OK")
		})

		c.Request, _ = http.NewRequest(http.MethodGet, tc.InputReqPath, nil)
		r.ServeHTTP(resp, c.Request)
	}
}

func TestRequestLogMiddlewareWithSlowRequest(t *testing.T) {
	t.Parallel()

	// Create a temporary directory for test traces
	tempDir, err := os.MkdirTemp("", "test-traces-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	lgr := logger.New("info", os.Stdout)
	fr := flightrecorder.New(lgr, tempDir, time.Second, 1<<20)
	require.NotNil(t, fr, "flight recorder should be created successfully")

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequestLogMiddleware(lgr, fr))

	// Add a slow endpoint
	router.GET("/slow", func(ctx *gin.Context) {
		time.Sleep(600 * time.Millisecond) // Exceeds SlowRequestThreshold
		ctx.String(http.StatusOK, "OK")
	})

	// Make request
	req, _ := http.NewRequest(http.MethodGet, "/slow", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	// Verify a trace file was created
	entries, err := os.ReadDir(tempDir)
	require.NoError(t, err)
	require.NotEmpty(t, entries, "at least one trace file should be created for slow request")
}
