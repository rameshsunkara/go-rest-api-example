package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestRequestLogMiddleware(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Create a mock logger
	mockLogger := logger.Setup("test")

	// Use the middleware with the mock logger
	router.Use(RequestLogMiddleware(mockLogger))

	// Define a test route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test")
	})

	// Create a test request to the defined route
	req, _ := http.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verify the middleware behavior
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.NotNil(t, mockLogger.reqID)
	assert.Equal(t, "GET", mockLogger.method)
	assert.Equal(t, "/test", mockLogger.url)
	assert.Equal(t, "/test", mockLogger.path)
	assert.Equal(t, "", mockLogger.userAgent) // Since we didn't set User-Agent in the request
	assert.Equal(t, http.StatusOK, mockLogger.respStatus)
	assert.True(t, mockLogger.elapsedMs > 0)
}
