package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/db/mocks"
	"github.com/rameshsunkara/go-rest-api-example/internal/handlers"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/data"
	"github.com/stretchr/testify/assert"
)

func TestNewSeedHandler(t *testing.T) {
	sd := handlers.NewDataSeedHandler(&mocks.MockOrdersDataService{})
	assert.IsType(t, &handlers.SeedHandler{}, sd)
}

func TestSeedDB_Success(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sd := handlers.NewDataSeedHandler(&mocks.MockOrdersDataService{
		CreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
			return "random-id", nil
		},
	})

	// Call actual function
	sd.SeedDB(c)

	resp := w.Result()

	// Check results
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
}

func TestSeedDB_Failure(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sd := handlers.NewDataSeedHandler(&mocks.MockOrdersDataService{
		CreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
			return "", assert.AnError
		},
	})

	// Call actual function
	sd.SeedDB(c)

	resp := w.Result()

	// Check results
	assert.EqualValues(t, http.StatusInternalServerError, resp.StatusCode)
}
