package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/mocks"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

var (
	svcInfo = &models.ServiceInfo{
		Name:        "test-api-service",
		Version:     "rams-fav",
		UpTime:      time.Now(),
		Environment: "test",
	}
)

func TestListOfRoutes(t *testing.T) {
	router := WebRouter(svcInfo, &mocks.MockDataMgr{})
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
		Path:   "/api/v1/orders",
	})

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodGet,
		Path:   "/api/v1/orders/:id",
	})

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodPost,
		Path:   "/api/v1/orders",
	})

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodPut,
		Path:   "/api/v1/orders",
	})

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodDelete,
		Path:   "/api/v1/orders/:id",
	})

}

func TestModeSpecificRoutes(t *testing.T) {
	svcInfo.Environment = "dev"
	router := WebRouter(svcInfo, &mocks.MockDataMgr{})
	list := router.Routes()
	mode := gin.Mode()

	assert.Equal(t, gin.DebugMode, mode)

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodPost,
		Path:   "/seedDB",
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
