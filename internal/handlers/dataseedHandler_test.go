package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/db/mocks"
	"github.com/rameshsunkara/go-rest-api-example/internal/handlers"
	"github.com/rameshsunkara/go-rest-api-example/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewSeedHandler(t *testing.T) {
	sd := handlers.NewSeedController(&mocks.MockOrdersDataService{})
	assert.IsType(t, &handlers.SeedController{}, sd)
	assert.IsType(t, &mocks.MockOrdersDataService{}, sd.OrdersDataSvc)
}

func TestSeedDB(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	mocks.CreateFunc = func(ctx context.Context, purchaseOrder *types.Order) (string, error) {
		return "random-id", nil
	}
	sd := handlers.NewSeedController(&mocks.MockOrdersDataService{})

	// Call actual function
	sd.SeedDB(c)

	resp := w.Result()

	// Check results
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
}
