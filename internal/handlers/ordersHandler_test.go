package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/db/mocks"
	errors2 "github.com/rameshsunkara/go-rest-api-example/internal/errors"
	"github.com/rameshsunkara/go-rest-api-example/internal/handlers"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/data"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/external"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func UnMarshalOrderData(d []byte) (*data.Order, error) {
	var r data.Order
	err := json.Unmarshal(d, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func UnMarshalOrdersData(d []byte) (*[]data.Order, error) {
	var r []data.Order
	err := json.Unmarshal(d, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func setupTestContext() (*gin.Context, *gin.Engine, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	return c, r, recorder
}

func TestOrdersHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		input          external.OrderInput
		mockCreateFunc func(context.Context, *data.Order) (string, error)
		expectedCode   int
		expectedError  *external.APIError
	}{
		{
			name: "Success",
			input: external.OrderInput{
				Products: []external.ProductInput{
					{Name: "Product 1", Price: 10.0, Quantity: 2},
				},
			},
			mockCreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
				return "1", nil
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "Invalid Input",
			input: external.OrderInput{
				Products: []external.ProductInput{
					{Name: "", Price: 10.0, Quantity: 2},
				},
			},
			mockCreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
				return "", errors.New("invalid input")
			},
			expectedCode: http.StatusInternalServerError, // use this once code is updated properly http.StatusBadRequest,
			expectedError: &external.APIError{
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorCode:      errors2.OrderCreateServerError,
				Message:        "unexpected Error occurred, please try again later",
			},
		},
		{
			name: "Internal Server Error",
			input: external.OrderInput{
				Products: []external.ProductInput{
					{Name: "Product 1", Price: 10.0, Quantity: 2},
				},
			},
			mockCreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
				return "", errors.New(errors2.UnexpectedErrorMessage)
			},
			expectedCode: http.StatusInternalServerError,
			expectedError: &external.APIError{
				HTTPStatusCode: http.StatusInternalServerError,
				Message:        errors2.UnexpectedErrorMessage,
				ErrorCode:      errors2.OrderCreateServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lgr := logger.Setup(models.ServiceEnv{Name: "test"})
			c, r, recorder := setupTestContext()
			handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
				CreateFunc: tt.mockCreateFunc,
			}, lgr)
			r.POST("/orders", handler.Create)

			body, _ := json.Marshal(tt.input)
			c.Request, _ = http.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
			r.ServeHTTP(recorder, c.Request)

			assert.Equal(t, tt.expectedCode, recorder.Code)

			if tt.expectedError != nil {
				var apiErr external.APIError
				err := json.Unmarshal(recorder.Body.Bytes(), &apiErr)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError.HTTPStatusCode, apiErr.HTTPStatusCode)
				assert.Equal(t, tt.expectedError.ErrorCode, apiErr.ErrorCode)
				assert.Equal(t, tt.expectedError.Message, apiErr.Message)
			} else {
				var responseOrder external.Order
				err := json.Unmarshal(recorder.Body.Bytes(), &responseOrder)
				require.NoError(t, err)
				assert.Equal(t, int64(1), responseOrder.Version)
				assert.NotNil(t, responseOrder.CreatedAt)
				assert.NotNil(t, responseOrder.UpdatedAt)
				assert.Equal(t, tt.input.Products[0].Name, responseOrder.Products[0].Name)
				assert.InEpsilon(t, tt.input.Products[0].Price, responseOrder.Products[0].Price, 0)
				assert.Equal(t, tt.input.Products[0].Quantity, responseOrder.Products[0].Quantity)
				assert.Equal(t, "1", responseOrder.ID)
				assert.InEpsilon(t, 20.0, responseOrder.TotalAmount, 0)
				assert.Equal(t, data.OrderPending, responseOrder.Status)
			}
		})
	}
}

func TestOrdersHandler_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		limit          string
		mockGetAllFunc func(context.Context, int64) (*[]data.Order, error)
		expectedCode   int
		expectedError  *external.APIError
		expectedLength int
	}{
		{
			name:  "Success",
			limit: "10",
			mockGetAllFunc: func(_ context.Context, _ int64) (*[]data.Order, error) {
				dataBytes, err := os.ReadFile("../mockData/orders.json")
				if err != nil {
					return nil, err
				}
				dataOrders, _ := UnMarshalOrdersData(dataBytes)
				return dataOrders, nil
			},
			expectedCode:   http.StatusOK,
			expectedLength: 10,
		},
		{
			name:  "DB Read Failure",
			limit: "10",
			mockGetAllFunc: func(_ context.Context, _ int64) (*[]data.Order, error) {
				return nil, errors.New("db error")
			},
			expectedCode: http.StatusInternalServerError,
			expectedError: &external.APIError{
				HTTPStatusCode: http.StatusInternalServerError,
				Message:        errors2.UnexpectedErrorMessage,
			},
		},
		{
			name:  "Limit Out of Bounds",
			limit: "10000",
			mockGetAllFunc: func(_ context.Context, _ int64) (*[]data.Order, error) {
				return nil, nil
			},
			expectedCode: http.StatusBadRequest,
			expectedError: &external.APIError{
				HTTPStatusCode: http.StatusBadRequest,
				Message:        "Integer value within 1 and 100 is expected for limit query param",
			},
		},
		{
			name:  "Invalid Limit",
			limit: "ABC",
			mockGetAllFunc: func(_ context.Context, _ int64) (*[]data.Order, error) {
				return nil, nil
			},
			expectedCode: http.StatusBadRequest,
			expectedError: &external.APIError{
				HTTPStatusCode: http.StatusBadRequest,
				Message:        "Integer value within 1 and 100 is expected for limit query param",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lgr := logger.Setup(models.ServiceEnv{Name: "test"})
			c, r, recorder := setupTestContext()
			handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
				GetAllFunc: tt.mockGetAllFunc,
			}, lgr)
			r.GET("/orders", handler.GetAll)

			c.Request, _ = http.NewRequest(http.MethodGet, "/orders", nil)
			q := c.Request.URL.Query()
			q.Add("limit", tt.limit)
			c.Request.URL.RawQuery = q.Encode()

			r.ServeHTTP(recorder, c.Request)

			assert.Equal(t, tt.expectedCode, recorder.Code)

			if tt.expectedError != nil {
				var apiErr external.APIError
				err := json.Unmarshal(recorder.Body.Bytes(), &apiErr)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError.HTTPStatusCode, apiErr.HTTPStatusCode)
				assert.Equal(t, tt.expectedError.Message, apiErr.Message)
			} else {
				var respOrders []external.Order
				err := json.Unmarshal(recorder.Body.Bytes(), &respOrders)
				require.NoError(t, err)
				assert.Len(t, respOrders, tt.expectedLength)
			}
		})
	}
}
