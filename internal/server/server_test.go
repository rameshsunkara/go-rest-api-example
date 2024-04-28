package server_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/db/mocks"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/rameshsunkara/go-rest-api-example/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestListOfRoutes(t *testing.T) {
	svcInfo := models.ServiceEnv{
		Name: "test",
		Port: "8080",
	}
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	router := server.WebRouter(svcInfo, &mocks.MockMongoMgr{}, lgr)
	list := router.Routes()
	mode := gin.Mode()

	assert.Equal(t, gin.ReleaseMode, mode)

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodGet,
		Path:   "/status",
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
	svcInfo := models.ServiceEnv{
		Name: "dev",
		Port: "8080",
	}
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	router := server.WebRouter(svcInfo, &mocks.MockMongoMgr{}, lgr)
	list := router.Routes()
	mode := gin.Mode()

	assert.Equal(t, gin.DebugMode, mode)

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodPost,
		Path:   "/internal/seed-local-db",
	})
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
