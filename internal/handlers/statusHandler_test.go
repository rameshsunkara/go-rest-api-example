package handlers_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/bogdanutanu/go-rest-api-example/internal/db/mocks"
	"github.com/bogdanutanu/go-rest-api-example/internal/handlers"
	"github.com/bogdanutanu/go-rest-api-example/pkg/logger"
	"github.com/bogdanutanu/go-rest-api-example/pkg/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStatusHandler(t *testing.T) {
	t.Parallel()
	mockMgr := &mocks.MockMongoMgr{}
	tests := []struct {
		name    string
		lgr     logger.Logger
		mgr     mongodb.MongoManager
		wantErr bool
	}{
		{
			name:    "success",
			lgr:     lgr,
			mgr:     mockMgr,
			wantErr: false,
		},
		{
			name:    "nil logger",
			lgr:     nil,
			mgr:     mockMgr,
			wantErr: true,
		},
		{
			name:    "nil manager",
			lgr:     lgr,
			mgr:     nil,
			wantErr: true,
		},
		{
			name:    "nil logger and manager",
			lgr:     nil,
			mgr:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h, err := handlers.NewStatusHandler(tt.lgr, tt.mgr)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, h)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, h)
			}
		})
	}
}

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
			expectedCode: http.StatusOK,
		},
		{
			name: "StatusDown",
			mockPingFunc: func() error {
				return errors.New("DB Connection Failed")
			},
			expectedCode: http.StatusOK,
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
