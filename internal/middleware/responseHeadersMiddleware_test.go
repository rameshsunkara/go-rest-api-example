package middleware_test

import (
	"net/http"
	"net/http/httptest"
	_ "strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestResponseHeadersMiddleware(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Use the middleware
	router.Use(middleware.ResponseHeadersMiddleware())

	// Define a test route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test")
	})

	// Create a test request to the defined route
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verify the response headers
	assert.Equal(t, "SAMEORIGIN", resp.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", resp.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "max-age=31536000; preload", resp.Header().Get("Strict-Transport-Security"),
		"All expected headers should be set")
}

func TestResponseHeadersMiddleware_CustomHeaders(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Use the middleware
	router.Use(middleware.ResponseHeadersMiddleware())

	// Define a test route
	router.GET("/test", func(c *gin.Context) {
		c.Writer.Header().Set("Custom-Header", "custom-value")
		c.String(http.StatusOK, "Test")
	})

	// Create a test request to the defined route
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verify the custom header and the standard headers
	assert.Equal(t, "SAMEORIGIN", resp.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", resp.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "max-age=31536000; preload", resp.Header().Get("Strict-Transport-Security"))
	assert.Equal(t, "custom-value", resp.Header().Get("Custom-Header"),
		"Custom-Header should be set along with Standard Headers")
}

func TestResponseHeadersMiddleware_NoCache(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Use the middleware
	router.Use(middleware.ResponseHeadersMiddleware())

	// Define a test route
	router.GET("/test", func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-store")
		c.String(http.StatusOK, "Test")
	})

	// Create a test request to the defined route
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verify the Cache-Control header should not be overwritten
	assert.Equal(t, "no-store", resp.Header().Get("Cache-Control"),
		"Cache-Control should not be overwritten")
}

func TestResponseHeadersMiddleware_NoHeadersSet(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Use the middleware
	router.Use(middleware.ResponseHeadersMiddleware())

	// Define a test route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test")
	})

	// Create a test request to the defined route
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verify the absence of custom headers
	assert.Empty(t, resp.Header().Get("Custom-Header"), "Custom-Header should not be set")
}
