package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/bogdanutanu/go-rest-api-example/internal/db"
	"github.com/bogdanutanu/go-rest-api-example/internal/db/mocks"
	"github.com/bogdanutanu/go-rest-api-example/internal/handlers"
	"github.com/bogdanutanu/go-rest-api-example/internal/models/data"
	"github.com/bogdanutanu/go-rest-api-example/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDataSeedHandler(t *testing.T) {
	t.Parallel()
	mockSvc := &mocks.MockOrdersDataService{}
	tests := []struct {
		name    string
		lgr     logger.Logger
		svc     db.OrdersDataService
		wantErr bool
	}{
		{
			name:    "success",
			lgr:     lgr,
			svc:     mockSvc,
			wantErr: false,
		},
		{
			name:    "nil logger",
			lgr:     nil,
			svc:     mockSvc,
			wantErr: true,
		},
		{
			name:    "nil service",
			lgr:     lgr,
			svc:     nil,
			wantErr: true,
		},
		{
			name:    "nil logger and service",
			lgr:     nil,
			svc:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h, err := handlers.NewDataSeedHandler(tt.lgr, tt.svc)
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

func TestDataSeedHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		mockCreateFunc func(context.Context, *data.Order) (string, error)
		expectedCode   int
	}{
		{
			name: "SeedDB_Success",
			mockCreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
				return "random-id", nil
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "SeedDB_Failure",
			mockCreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
				return "", errors.New("create error")
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // mark the test as capable of running in parallel

			c, r, recorder := setupTestContext()
			sd, err := handlers.NewDataSeedHandler(lgr, &mocks.MockOrdersDataService{
				CreateFunc: tt.mockCreateFunc,
			})
			if err != nil {
				t.Errorf("failed to create dataseed handler")
				return
			}
			r.POST("/seed", sd.SeedDB)

			c.Request, _ = http.NewRequest(http.MethodPost, "/seed", nil)
			r.ServeHTTP(recorder, c.Request)

			resp := recorder.Result()
			require.NoError(t, resp.Body.Close())
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}
