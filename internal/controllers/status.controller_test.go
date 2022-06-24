package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/mocks"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func UnMarshalStatusResponse(resp *http.Response) (StatusResponse, error) {
	body, _ := io.ReadAll(resp.Body)
	var statusResponse StatusResponse
	err := json.Unmarshal(body, &statusResponse)
	return statusResponse, err
}

var (
	svcInfo = &models.ServiceInfo{
		Name:        "test-api-service",
		Version:     "rams-fav",
		UpTime:      time.Now(),
		Environment: "test",
	}
	s = NewStatusController(svcInfo, &mocks.MockDataMgr{})
)

func TestStatusSuccess(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	mocks.PingFunc = func() error {
		return nil
	}

	// Call actual function
	s.CheckStatus(c)

	// Check results
	resp := w.Result()
	statusResponse, err := UnMarshalStatusResponse(resp)
	if err != nil {
		t.Fail()
	}
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
	assert.EqualValues(t, "test", statusResponse.Environment)
}

func TestStatusDown(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mocks.PingFunc = func() error {
		return errors.New("DB Connection Failed")
	}

	s.CheckStatus(c)

	resp := w.Result()
	statusResponse, err := UnMarshalStatusResponse(resp)
	if err != nil {
		t.Fail()
	}

	assert.EqualValues(t, http.StatusFailedDependency, resp.StatusCode)
	assert.EqualValues(t, "rams-fav", statusResponse.Version)
}
