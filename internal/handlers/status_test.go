package handlers_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/db/mocks"
	"github.com/rameshsunkara/go-rest-api-example/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func UnMarshalStatusResponse(resp *http.Response) (string, error) {
	body, _ := io.ReadAll(resp.Body)
	var statusResponse string
	err := json.Unmarshal(body, &statusResponse)
	return statusResponse, err
}

func TestStatusSuccess(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	mocks.PingFunc = func() error {
		return nil
	}
	s := handlers.NewStatusController(&mocks.MockMongoMgr{})

	// Call actual function
	s.CheckStatus(c)

	// Parse results
	resp := w.Result()
	statusResponse, err := UnMarshalStatusResponse(resp)
	if err != nil {
		t.Fail()
	}
	// Check results
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
	assert.EqualValues(t, handlers.UP, statusResponse)
}

func TestStatusDown(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	mocks.PingFunc = func() error {
		return errors.New("DB Connection Failed")
	}
	s := handlers.NewStatusController(&mocks.MockMongoMgr{})

	// Call actual function
	s.CheckStatus(c)

	// Parse results
	resp := w.Result()
	statusResponse, err := UnMarshalStatusResponse(resp)
	if err != nil {
		t.Fail()
	}
	// Check results
	assert.EqualValues(t, http.StatusFailedDependency, resp.StatusCode)
	assert.EqualValues(t, handlers.DOWN, statusResponse)
}
