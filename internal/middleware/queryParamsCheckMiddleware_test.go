package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/middleware"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestQueryParamsCheckMiddleware_ValidParams(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(resp)
	r.Use(middleware.QueryParamsCheckMiddleware(lgr))

	reqURL := "/ecommerce/v1/orders"
	r.GET(reqURL, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	c.Request, _ = http.NewRequest(http.MethodGet, reqURL, nil)
	q := c.Request.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "5")
	c.Request.URL.RawQuery = q.Encode()
	r.ServeHTTP(resp, c.Request)

	// Assert response
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestQueryParamsCheckMiddleware_InvalidParams(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(resp)
	r.Use(middleware.QueryParamsCheckMiddleware(lgr))

	reqURL := "/ecommerce/v1/orders"
	r.GET(reqURL, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	c.Request, _ = http.NewRequest(http.MethodGet, reqURL, nil)
	q := c.Request.URL.Query()
	q.Add("example", "10")
	c.Request.URL.RawQuery = q.Encode()
	r.ServeHTTP(resp, c.Request)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestQueryParamsCheckMiddleware_UnregisteredPath(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(resp)
	r.Use(middleware.QueryParamsCheckMiddleware(lgr))

	reqURL := "/ecommerce/v2/orders"
	r.GET(reqURL, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	c.Request, _ = http.NewRequest(http.MethodGet, reqURL, nil)
	q := c.Request.URL.Query()
	q.Add("example", "10")
	c.Request.URL.RawQuery = q.Encode()
	r.ServeHTTP(resp, c.Request)

	// Assert response
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestHasUnSupportedQueryParams(t *testing.T) {
	testCases := []struct {
		description     string
		queryParams     url.Values
		supportedParams map[string]bool
		expectedVal     bool
	}{
		{
			description:     "All parameters are supported",
			queryParams:     url.Values{"param1": []string{"value1"}, "param2": []string{"value2"}},
			supportedParams: map[string]bool{"param1": true, "param2": true},
			expectedVal:     false,
		},
		{
			description:     "Some parameters are not supported",
			queryParams:     url.Values{"param1": []string{"value1"}, "param3": []string{"value3"}},
			supportedParams: map[string]bool{"param1": true, "param2": true},
			expectedVal:     true,
		},
		{
			description:     "No parameters are supported",
			queryParams:     url.Values{"param1": []string{"value1"}, "param3": []string{"value3"}},
			supportedParams: map[string]bool{},
			expectedVal:     true,
		},
		{
			description:     "handle when nil is passed as supportedParams",
			queryParams:     url.Values{"param1": []string{"value1"}, "param3": []string{"value3"}},
			supportedParams: nil,
			expectedVal:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			req := &http.Request{URL: &url.URL{RawQuery: tc.queryParams.Encode()}}
			supported := middleware.HasUnSupportedQueryParams(req, tc.supportedParams)
			if supported != tc.expectedVal {
				t.Errorf("Expected %v, but got %v", tc.expectedVal, supported)
			}
		})
	}
}
