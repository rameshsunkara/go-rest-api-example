package handlers_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/rameshsunkara/go-rest-api-example/internal/db/mocks"
	"github.com/rameshsunkara/go-rest-api-example/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		mockPingFunc func() error
		expectedCode int
	}{
		{
			name: "StatusSuccess",
			mockPingFunc: func() error {
				return nil
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name: "StatusDown",
			mockPingFunc: func() error {
				return errors.New("DB Connection Failed")
			},
			expectedCode: http.StatusFailedDependency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // mark the test as capable of running in parallel

			// Test Setup
			c, _, recorder := setupTestContext()
			s, err := handlers.NewStatusHandler(lgr, &mocks.MockMongoMgr{
				PingFunc: tt.mockPingFunc,
			})

			// Call actual function
			s.CheckStatus(c)

			// Parse results
			resp := recorder.Result()

			// Check results
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
			require.NoError(t, err)
		})
	}
}
