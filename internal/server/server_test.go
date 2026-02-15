package server_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/bogdanutanu/go-rest-api-example/internal/config"
	"github.com/bogdanutanu/go-rest-api-example/internal/db/mocks"
	"github.com/bogdanutanu/go-rest-api-example/internal/server"
	"github.com/bogdanutanu/go-rest-api-example/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListOfRoutes(t *testing.T) {
	svcInfo := &config.ServiceEnvConfig{
		Environment:          "test",
		Port:                 "8080",
		LogLevel:             "info",
		DBCredentialsSideCar: "/path/to/mongo/sidecar",
		DBHosts:              "localhost",
		DBName:               "testDB",
	}
	lgr := logger.New("info", os.Stdout)
	router, err := server.WebRouter(svcInfo, lgr, &mocks.MockMongoMgr{})
	if err != nil {
		t.Errorf("failed to create WebRouter")
		return
	}
	list := router.Routes()
	mode := gin.Mode()

	assert.Equal(t, gin.ReleaseMode, mode)

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodGet,
		Path:   "/healthz",
	})

	assertRouteNotPresent(t, list, gin.RouteInfo{
		Method: http.MethodPost,
		Path:   "/seedDB",
	})

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodGet,
		Path:   "/ecommerce/v1/orders",
	})

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodGet,
		Path:   "/ecommerce/v1/orders/:id",
	})

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodPost,
		Path:   "/ecommerce/v1/orders",
	})

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodDelete,
		Path:   "/ecommerce/v1/orders/:id",
	})
}

func TestModeSpecificRoutes(t *testing.T) {
	svcInfo := &config.ServiceEnvConfig{
		Environment: "dev",
		Port:        "8080",
	}
	lgr := logger.New("info", os.Stdout)
	router, err := server.WebRouter(svcInfo, lgr, &mocks.MockMongoMgr{})
	if err != nil {
		t.Errorf("failed to create WebRouter")
		return
	}
	list := router.Routes()
	mode := gin.Mode()

	assert.Equal(t, gin.DebugMode, mode)

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodPost,
		Path:   "/internal/seed-local-db",
	})
}

func TestWebRouterWithTracingEnabled(t *testing.T) {
	svcInfo := &config.ServiceEnvConfig{
		Environment:          "test",
		Port:                 "8080",
		LogLevel:             "info",
		DBCredentialsSideCar: "/path/to/mongo/sidecar",
		DBHosts:              "localhost",
		DBName:               "testDB",
		EnableTracing:        true, // Enable tracing
	}
	lgr := logger.New("info", os.Stdout)
	router, err := server.WebRouter(svcInfo, lgr, &mocks.MockMongoMgr{})

	require.NoError(t, err)
	assert.NotNil(t, router)

	// Verify router is properly configured
	list := router.Routes()
	assert.NotEmpty(t, list)
}

func assertRoutePresent(t *testing.T, gotRoutes gin.RoutesInfo, wantRoute gin.RouteInfo) {
	for _, gotRoute := range gotRoutes {
		if gotRoute.Path == wantRoute.Path && gotRoute.Method == wantRoute.Method {
			return
		}
	}
	t.Errorf("route not found: %v", wantRoute)
}

func assertRouteNotPresent(t *testing.T, gotRoutes gin.RoutesInfo, wantRoute gin.RouteInfo) {
	for _, gotRoute := range gotRoutes {
		if gotRoute.Path == wantRoute.Path && gotRoute.Method == wantRoute.Method {
			t.Errorf("route found: %v", wantRoute)
		}
	}
}
