package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/db/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	createFunc func(order interface{}) (*mongo.InsertOneResult, error)
	sd         = NewSeedController(&mocks.MockOrdersDataService{})
)

func TestNewSeedHandler(t *testing.T) {
	assert.IsType(t, &SeedController{}, sd)
	assert.IsType(t, &mocks.MockOrdersDataService{}, sd.dataSvc)
}

func TestSeedDB(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	mocks.CreateFunc = func(ctx context.Context, purchaseOrder interface{}) (*mongo.InsertOneResult, error) {
		return nil, nil
	}

	// Call actual function
	sd.SeedDB(c)

	resp := w.Result()

	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
}
