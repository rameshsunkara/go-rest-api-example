package handlers_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/rameshsunkara/go-rest-api-example/internal/db/mocks"
	"github.com/rameshsunkara/go-rest-api-example/internal/handlers"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			sd := handlers.NewDataSeedHandler(&mocks.MockOrdersDataService{
				CreateFunc: tt.mockCreateFunc,
			})
			r.POST("/seed", sd.SeedDB)

			c.Request, _ = http.NewRequest(http.MethodPost, "/seed", nil)
			r.ServeHTTP(recorder, c.Request)

			resp := recorder.Result()
			require.NoError(t, resp.Body.Close())
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}
