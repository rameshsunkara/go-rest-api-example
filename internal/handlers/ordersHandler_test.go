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
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func TestOrdersHandler_Create_Success(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		CreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
			return "1", nil
		},
	}, lgr)
	r.POST("/orders", handler.Create)

	orderInput := external.OrderInput{
		Products: []external.ProductInput{
			{Name: "Product 1", Price: 10.0, Quantity: 2},
		},
	}

	body, _ := json.Marshal(orderInput)
	c.Request, _ = http.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	r.ServeHTTP(recorder, c.Request)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var responseOrder external.Order
	err := json.Unmarshal(recorder.Body.Bytes(), &responseOrder)
	require.NoError(t, err)

	assert.Equal(t, int64(1), responseOrder.Version)
	assert.NotNil(t, responseOrder.CreatedAt)
	assert.NotNil(t, responseOrder.UpdatedAt)
	assert.Equal(t, orderInput.Products[0].Name, responseOrder.Products[0].Name)
	assert.InEpsilon(t, orderInput.Products[0].Price, responseOrder.Products[0].Price, 0)
	assert.Equal(t, orderInput.Products[0].Quantity, responseOrder.Products[0].Quantity)
	assert.Equal(t, "1", responseOrder.ID)
	assert.InEpsilon(t, 20.0, responseOrder.TotalAmount, 0)
	assert.Equal(t, data.OrderPending, responseOrder.Status)
}

func TestOrdersHandler_Create_InvalidInput(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		CreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
			return "mock_order_id", nil
		},
	}, lgr)
	r.POST("/orders", handler.Create)
	invalidInput := "{ invalid JSON }"
	c.Request, _ = http.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte(invalidInput)))

	r.ServeHTTP(recorder, c.Request)

	var apiErr external.APIError
	err := json.Unmarshal(recorder.Body.Bytes(), &apiErr)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatusCode)
	assert.Equal(t, "orders_create_invalid_input", apiErr.ErrorCode)
	assert.Equal(t, "Invalid order request body", apiErr.Message)
}

func TestOrdersHandler_Create_InternalServerError(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		CreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
			return "", errors.New(errors2.UnexpectedErrorMessage)
		},
	}, lgr)
	r.POST("/orders", handler.Create)
	orderInput := external.OrderInput{
		Products: []external.ProductInput{
			{Name: "Product 1", Price: 10.0, Quantity: 2},
		},
	}
	body, _ := json.Marshal(orderInput)
	c.Request, _ = http.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))

	r.ServeHTTP(recorder, c.Request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	var apiErr external.APIError
	err := json.Unmarshal(recorder.Body.Bytes(), &apiErr)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, apiErr.HTTPStatusCode)
	assert.Equal(t, errors2.UnexpectedErrorMessage, apiErr.Message)
}

func TestGetAllOrdersSuccess(t *testing.T) {
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		GetAllFunc: func(_ context.Context, _ int64) (*[]data.Order, error) {
			dataBytes, err := os.ReadFile("../mockData/orders.json")
			if err != nil {
				return nil, err
			}
			dataOrders, _ := UnMarshalOrdersData(dataBytes)
			return dataOrders, nil
		},
	}, lgr)
	r.GET("/orders", handler.GetAll)
	c.Request, _ = http.NewRequest(http.MethodGet, "/orders", nil)

	r.ServeHTTP(recorder, c.Request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	var respOrders []external.Order
	err := json.Unmarshal(recorder.Body.Bytes(), &respOrders)
	require.NoError(t, err)
	assert.Len(t, respOrders, 10)
}

func TestGetAllOrdersFailure_DBRead(t *testing.T) {
	// Test Setup
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		GetAllFunc: func(_ context.Context, _ int64) (*[]data.Order, error) {
			dataBytes, err := os.ReadFile("../mockData/non-existent.json")
			if err != nil {
				return nil, err
			}
			dataOrders, _ := UnMarshalOrdersData(dataBytes)
			return dataOrders, nil
		},
	}, lgr)
	r.GET("/orders", handler.GetAll)
	c.Request, _ = http.NewRequest(http.MethodGet, "/orders", nil)

	r.ServeHTTP(recorder, c.Request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestGetAllOrdersFailure_LimitOutOfBounds(t *testing.T) {
	// Test Setup
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{}, lgr)
	r.GET("/orders", handler.GetAll)
	c.Request, _ = http.NewRequest(http.MethodGet, "/orders", nil)
	q := c.Request.URL.Query()
	q.Add("limit", "10000")
	c.Request.URL.RawQuery = q.Encode()

	r.ServeHTTP(recorder, c.Request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetAllOrdersFailure_InvalidLimit(t *testing.T) {
	// Test Setup
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	gin.SetMode(gin.TestMode)
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{}, lgr)
	r.GET("/orders", handler.GetAll)
	c.Request, _ = http.NewRequest(http.MethodGet, "/orders", nil)
	q := c.Request.URL.Query()
	q.Add("limit", "ABC")
	c.Request.URL.RawQuery = q.Encode()

	r.ServeHTTP(recorder, c.Request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetOrderByIDSuccess(t *testing.T) {
	// Test Setup
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		GetByIDFunc: func(_ context.Context, _ primitive.ObjectID) (*data.Order, error) {
			dataBytes, err := os.ReadFile("../mockData/order.json")
			if err != nil {
				return nil, err
			}
			dataOrder, _ := UnMarshalOrderData(dataBytes)
			return dataOrder, nil
		},
	}, lgr)
	r.GET("/ecommerce/v1/orders/:id", handler.GetByID)
	c.Request, _ = http.NewRequest(http.MethodGet, "/ecommerce/v1/orders/609d9ed771df2a0d99bf0077", nil)

	r.ServeHTTP(recorder, c.Request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	var respOrder external.Order
	err := json.Unmarshal(recorder.Body.Bytes(), &respOrder)
	require.NoError(t, err)
	assert.Equal(t, "609d9ed771df2a0d99bf0077", respOrder.ID)
}

func TestGetOrderByID_DBReadFailure(t *testing.T) {
	// Test Setup
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		GetByIDFunc: func(_ context.Context, _ primitive.ObjectID) (*data.Order, error) {
			return nil, errors.New("db error")
		},
	}, lgr)
	r.GET("/ecommerce/v1/orders/:id", handler.GetByID)
	c.Request, _ = http.NewRequest(http.MethodGet, "/ecommerce/v1/orders/609d9ed771df2a0d99bf0077", nil)

	r.ServeHTTP(recorder, c.Request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestGetOrderByID_BadPathParam(t *testing.T) {
	// Test Setup
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		GetByIDFunc: func(_ context.Context, _ primitive.ObjectID) (*data.Order, error) {
			return nil, errors.New("db error")
		},
	}, lgr)
	r.GET("/ecommerce/v1/orders/:id", handler.GetByID)
	c.Request, _ = http.NewRequest(http.MethodGet, "/ecommerce/v1/orders/''", nil)

	r.ServeHTTP(recorder, c.Request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestDeleteOrderByIDSuccess(t *testing.T) {
	// Test Setup
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		DeleteByIDFunc: func(_ context.Context, _ primitive.ObjectID) error {
			return nil
		},
	}, lgr)
	r.DELETE("/ecommerce/v1/orders/:id", handler.DeleteByID)
	c.Request, _ = http.NewRequest(http.MethodDelete, "/ecommerce/v1/orders/609d9ed771df2a0d99bf0077", nil)

	r.ServeHTTP(recorder, c.Request)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestDeleteOrderByID_DBFailure(t *testing.T) {
	// Test Setup
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		DeleteByIDFunc: func(_ context.Context, _ primitive.ObjectID) error {
			return errors.New("db error")
		},
	}, lgr)
	r.DELETE("/ecommerce/v1/orders/:id", handler.DeleteByID)
	c.Request, _ = http.NewRequest(http.MethodDelete, "/ecommerce/v1/orders/609d9ed771df2a0d99bf0077", nil)

	r.ServeHTTP(recorder, c.Request)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestDeleteOrderByID_BadPathParam(t *testing.T) {
	// Test Setup
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	recorder := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(recorder)
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		DeleteByIDFunc: func(_ context.Context, _ primitive.ObjectID) error {
			return nil
		},
	}, lgr)
	r.DELETE("/ecommerce/v1/orders/:id", handler.DeleteByID)
	c.Request, _ = http.NewRequest(http.MethodDelete, "/ecommerce/v1/orders/''", nil)

	r.ServeHTTP(recorder, c.Request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
