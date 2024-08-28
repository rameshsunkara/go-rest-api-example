package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Use the middleware
	router.Use(middleware.AuthMiddleware())

	// Define a test route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test")
	})

	// Create a test request to the defined route
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verify that the response status code is OK
	assert.Equal(t, http.StatusOK, resp.Code, "Response status code should be OK")
}

func TestAuthMiddleware_WithNext(t *testing.T) {
	// Create a new Gin router
	router := gin.New()

	// Variable to check if Next() is called
	var nextCalled bool

	// Use the middleware
	router.Use(middleware.AuthMiddleware())

	// Define a test route
	router.GET("/test", func(c *gin.Context) {
		nextCalled = true
		c.String(http.StatusOK, "Test")
	})

	// Create a test request to the defined route
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Verify that the response status code is OK
	assert.Equal(t, http.StatusOK, resp.Code)

	// Ensure that Next() was called
	assert.True(t, nextCalled, "Next() should be called")
}
