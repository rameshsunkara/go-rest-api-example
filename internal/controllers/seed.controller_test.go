package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	createFunc func(order interface{}) (*mongo.InsertOneResult, error)
	sd         = NewSeedController(&mocks.MockDataService{})
)

func TestNewSeedHandler(t *testing.T) {
	assert.IsType(t, &SeedController{}, sd)
	assert.IsType(t, &mocks.MockDataService{}, sd.dataSvc)
}

func TestSeedDB(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	mocks.CreateFunc = func(purchaseOrder interface{}) (*mongo.InsertOneResult, error) {
		return nil, nil
	}

	// Call actual function
	sd.SeedDB(c)

	resp := w.Result()

	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
}
