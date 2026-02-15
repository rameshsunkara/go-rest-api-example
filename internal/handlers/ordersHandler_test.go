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

	"github.com/bogdanutanu/go-rest-api-example/internal/db"
	"github.com/bogdanutanu/go-rest-api-example/internal/db/mocks"
	errors2 "github.com/bogdanutanu/go-rest-api-example/internal/errors"
	"github.com/bogdanutanu/go-rest-api-example/internal/handlers"
	"github.com/bogdanutanu/go-rest-api-example/internal/models/data"
	"github.com/bogdanutanu/go-rest-api-example/internal/models/external"
	"github.com/bogdanutanu/go-rest-api-example/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var lgr logger.Logger

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	lgr = logger.New("debug", os.Stdout)
	os.Exit(m.Run())
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
	c, r := gin.CreateTestContext(recorder)
	return c, r, recorder
}

func TestNewOrdersHandler(t *testing.T) {
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
			h, err := handlers.NewOrdersHandler(tt.lgr, tt.svc)
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

func TestOrdersHandler_Create(t *testing.T) {
	t.Parallel()
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
			name:  "Invalid Input 2",
			input: external.OrderInput{},
			mockCreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
				return "", errors.New("invalid input")
			},
			expectedCode: http.StatusBadRequest, // use this once code is updated properly http.StatusBadRequest,
			expectedError: &external.APIError{
				HTTPStatusCode: http.StatusBadRequest,
				ErrorCode:      errors2.OrderCreateInvalidInput,
				Message:        "Invalid order request body",
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
			t.Parallel() // mark the test as capable of running in parallel

			c, r, recorder := setupTestContext()
			handler, err := handlers.NewOrdersHandler(lgr, &mocks.MockOrdersDataService{
				CreateFunc: tt.mockCreateFunc,
			})
			if err != nil {
				t.Errorf("failed to create orders handler")
				return
			}
			r.POST("/orders", handler.Create)

			body, _ := json.Marshal(tt.input)
			c.Request, _ = http.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
			r.ServeHTTP(recorder, c.Request)

			assert.Equal(t, tt.expectedCode, recorder.Code)

			if tt.expectedError != nil {
				var apiErr external.APIError
				respBodyErr := json.Unmarshal(recorder.Body.Bytes(), &apiErr)
				require.NoError(t, respBodyErr)
				assert.Equal(t, tt.expectedError.HTTPStatusCode, apiErr.HTTPStatusCode)
				assert.Equal(t, tt.expectedError.ErrorCode, apiErr.ErrorCode)
				assert.Equal(t, tt.expectedError.Message, apiErr.Message)
			} else {
				var responseOrder external.Order
				respBodyErr := json.Unmarshal(recorder.Body.Bytes(), &responseOrder)
				require.NoError(t, respBodyErr)
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
	t.Parallel()
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
				results := make([]data.Order, 10)
				return &results, nil
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
				results := make([]data.Order, 10)
				return &results, nil
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
			t.Parallel() // mark the test as capable of running in parallel

			c, r, recorder := setupTestContext()
			handler, err := handlers.NewOrdersHandler(lgr, &mocks.MockOrdersDataService{
				GetAllFunc: tt.mockGetAllFunc,
			})
			if err != nil {
				t.Errorf("failed to create orders handler")
				return
			}
			r.GET("/orders", handler.GetAll)

			c.Request, _ = http.NewRequest(http.MethodGet, "/orders", nil)
			q := c.Request.URL.Query()
			q.Add("limit", tt.limit)
			c.Request.URL.RawQuery = q.Encode()

			r.ServeHTTP(recorder, c.Request)

			assert.Equal(t, tt.expectedCode, recorder.Code)

			if tt.expectedError != nil {
				var apiErr external.APIError
				respBodyErr := json.Unmarshal(recorder.Body.Bytes(), &apiErr)
				require.NoError(t, respBodyErr)
				assert.Equal(t, tt.expectedError.HTTPStatusCode, apiErr.HTTPStatusCode)
				assert.Equal(t, tt.expectedError.Message, apiErr.Message)
			} else {
				var respOrders []external.Order
				respBodyErr := json.Unmarshal(recorder.Body.Bytes(), &respOrders)
				require.NoError(t, respBodyErr)
				assert.Len(t, respOrders, tt.expectedLength)
			}
		})
	}
}

func TestOrdersHandler_GetByID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		orderID         string
		mockGetByIDFunc func(ctx context.Context, oID primitive.ObjectID) (*data.Order, error)
		expectedCode    int
		expectedError   *external.APIError
	}{
		{
			name:    "Success",
			orderID: primitive.NewObjectID().Hex(),
			mockGetByIDFunc: func(_ context.Context, oID primitive.ObjectID) (*data.Order, error) {
				return &data.Order{ID: oID, Status: data.OrderPending}, nil
			},
			expectedCode: http.StatusOK,
		},
		{
			name:    "Not Found",
			orderID: primitive.NewObjectID().Hex(),
			mockGetByIDFunc: func(_ context.Context, _ primitive.ObjectID) (*data.Order, error) {
				return nil, errors.New("not found")
			},
			expectedCode: http.StatusInternalServerError,
			expectedError: &external.APIError{
				HTTPStatusCode: http.StatusInternalServerError,
				Message:        "failed to fetch order",
			},
		},
		{
			name:    "Zero Order Cannot be fetched",
			orderID: primitive.NilObjectID.Hex(),
			mockGetByIDFunc: func(_ context.Context, _ primitive.ObjectID) (*data.Order, error) {
				return nil, errors.New("not found")
			},
			expectedCode: http.StatusBadRequest,
			expectedError: &external.APIError{
				HTTPStatusCode: http.StatusBadRequest,
				Message:        "invalid order ID",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, r, recorder := setupTestContext()
			handler, err := handlers.NewOrdersHandler(lgr, &mocks.MockOrdersDataService{
				GetByIDFunc: tt.mockGetByIDFunc,
			})
			if err != nil {
				t.Errorf("failed to create orders handler")
				return
			}
			r.GET("/orders/:id", handler.GetByID)

			c.Request, _ = http.NewRequest(http.MethodGet, "/orders/"+tt.orderID, nil)
			r.ServeHTTP(recorder, c.Request)

			assert.Equal(t, tt.expectedCode, recorder.Code)

			if tt.expectedError != nil {
				var apiErr external.APIError
				respBodyErr := json.Unmarshal(recorder.Body.Bytes(), &apiErr)
				require.NoError(t, respBodyErr)
				assert.Equal(t, tt.expectedError.HTTPStatusCode, apiErr.HTTPStatusCode)
				assert.Equal(t, tt.expectedError.Message, apiErr.Message)
			} else {
				var responseOrder external.Order
				respBodyErr := json.Unmarshal(recorder.Body.Bytes(), &responseOrder)
				require.NoError(t, respBodyErr)
				assert.Equal(t, tt.orderID, responseOrder.ID)
			}
		})
	}
}

func TestOrdersHandler_Delete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		orderID        string
		mockDeleteFunc func(ctx context.Context, id primitive.ObjectID) error
		expectedCode   int
		expectedError  *external.APIError
	}{
		{
			name:    "Success",
			orderID: primitive.NewObjectID().Hex(),
			mockDeleteFunc: func(_ context.Context, _ primitive.ObjectID) error {
				return nil
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:    "Not Found",
			orderID: primitive.NewObjectID().Hex(),
			mockDeleteFunc: func(_ context.Context, _ primitive.ObjectID) error {
				return errors.New("not found")
			},
			expectedCode: http.StatusInternalServerError,
			expectedError: &external.APIError{
				HTTPStatusCode: http.StatusInternalServerError,
				Message:        "could not delete order",
			},
		},
		{
			name:    "Zero Order Cannot be Deleted",
			orderID: primitive.NilObjectID.Hex(),
			mockDeleteFunc: func(_ context.Context, _ primitive.ObjectID) error {
				return errors.New("not found")
			},
			expectedCode: http.StatusBadRequest,
			expectedError: &external.APIError{
				HTTPStatusCode: http.StatusBadRequest,
				Message:        "invalid order ID",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, r, recorder := setupTestContext()
			handler, err := handlers.NewOrdersHandler(lgr, &mocks.MockOrdersDataService{
				DeleteByIDFunc: tt.mockDeleteFunc,
			})
			if err != nil {
				t.Errorf("failed to create orders handler")
				return
			}
			r.DELETE("/orders/:id", handler.DeleteByID)

			c.Request, _ = http.NewRequest(http.MethodDelete, "/orders/"+tt.orderID, nil)
			r.ServeHTTP(recorder, c.Request)

			assert.Equal(t, tt.expectedCode, recorder.Code)

			if tt.expectedError != nil {
				var apiErr external.APIError
				respBodyErr := json.Unmarshal(recorder.Body.Bytes(), &apiErr)
				require.NoError(t, respBodyErr)
				assert.Equal(t, tt.expectedError.HTTPStatusCode, apiErr.HTTPStatusCode)
				assert.Equal(t, tt.expectedError.Message, apiErr.Message)
			}
		})
	}
}
