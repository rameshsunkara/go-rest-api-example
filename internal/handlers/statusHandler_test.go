package handlers_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/rameshsunkara/go-rest-api-example/internal/db/mocks"
	"github.com/rameshsunkara/go-rest-api-example/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func UnMarshalStatusResponse(resp *http.Response) (handlers.ServiceStatus, error) {
	body, _ := io.ReadAll(resp.Body)
	var statusResponse handlers.ServiceStatus
	err := json.Unmarshal(body, &statusResponse)
	return statusResponse, err
}

func TestStatusHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		mockPingFunc   func() error
		expectedStatus handlers.ServiceStatus
		expectedCode   int
	}{
		{
			name: "StatusSuccess",
			mockPingFunc: func() error {
				return nil
			},
			expectedStatus: handlers.UP,
			expectedCode:   http.StatusOK,
		},
		{
			name: "StatusDown",
			mockPingFunc: func() error {
				return errors.New("DB Connection Failed")
			},
			expectedStatus: handlers.DOWN,
			expectedCode:   http.StatusFailedDependency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // mark the test as capable of running in parallel

			// Test Setup
			c, _, recorder := setupTestContext()
			mocks.PingFunc = tt.mockPingFunc
			s := handlers.NewStatusController(&mocks.MockMongoMgr{})

			// Call actual function
			s.CheckStatus(c)

			// Parse results
			resp := recorder.Result()
			statusResponse, err := UnMarshalStatusResponse(resp)
			require.NoError(t, err)

			// Check results
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
			assert.Equal(t, tt.expectedStatus, statusResponse)
		})
	}
}
