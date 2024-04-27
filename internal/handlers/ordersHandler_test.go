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

func UnMarshalOrderResponse(d []byte) (*data.Order, error) {
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
	gin.SetMode(gin.TestMode)
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	r := gin.New()
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
	req, err := http.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var responseOrder external.Order
	err = json.Unmarshal(w.Body.Bytes(), &responseOrder)
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
	gin.SetMode(gin.TestMode)
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	r := gin.New()
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{
		CreateFunc: func(_ context.Context, _ *data.Order) (string, error) {
			return "mock_order_id", nil
		},
	}, lgr)
	r.POST("/orders", handler.Create)

	invalidInput := "{ invalid JSON }"
	req, err := http.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte(invalidInput)))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var apiErr external.APIError
	err = json.Unmarshal(w.Body.Bytes(), &apiErr)
	require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, apiErr.HTTPStatusCode)
	assert.Equal(t, "orders_create_invalid_input", apiErr.ErrorCode)
	assert.Equal(t, "Invalid order request body", apiErr.Message)
}

func TestOrdersHandler_Create_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	r := gin.New()
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
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var apiErr external.APIError
	err := json.Unmarshal(w.Body.Bytes(), &apiErr)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, apiErr.HTTPStatusCode)
	assert.Equal(t, "orders_create_server_error", apiErr.ErrorCode)
	assert.Equal(t, errors2.UnexpectedErrorMessage, apiErr.Message) // Assuming errors.UnexpectedErrorMessage
}

func TestGetAllOrdersSuccess(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	r := gin.New()
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

	req, err := http.NewRequest(http.MethodGet, "/orders", nil)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var respOrders []external.Order
	err = json.Unmarshal(w.Body.Bytes(), &respOrders)
	require.NoError(t, err)
	assert.Len(t, respOrders, 10)
}

func TestGetAllOrdersFailure_DBRead(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	r := gin.New()
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

	req, err := http.NewRequest(http.MethodGet, "/orders", nil)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetAllOrdersFailure_LimitOutOfBounds(t *testing.T) {
	// Test Setup
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(resp)
	gin.SetMode(gin.TestMode)
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{}, lgr)
	r.GET("/orders", handler.GetAll)

	c.Request, _ = http.NewRequest(http.MethodGet, "/orders", nil)
	q := c.Request.URL.Query()
	q.Add("limit", "10000")
	c.Request.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()
	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAllOrdersFailure_InvalidLimit(t *testing.T) {
	// Test Setup
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(resp)
	gin.SetMode(gin.TestMode)
	lgr := logger.Setup(models.ServiceEnv{Name: "test"})
	handler := handlers.NewOrdersHandler(&mocks.MockOrdersDataService{}, lgr)
	r.GET("/orders", handler.GetAll)

	c.Request, _ = http.NewRequest(http.MethodGet, "/orders", nil)
	q := c.Request.URL.Query()
	q.Add("limit", "ABC")
	c.Request.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()
	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

/*

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/db/mocks"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UnMarshalOrdersResponse(d []byte) (*[]models.Order, error) {
	var orders *[]models.Order
	err := json.Unmarshal(d, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func UnMarshalOrderResponse(d []byte) (*models.Order, error) {
	var r *models.Order
	err := json.Unmarshal(d, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func UnMarshalCreateOrderResponse(d []byte) (*mongo.InsertOneResult, error) {
	var r *mongo.InsertOneResult
	err := json.Unmarshal(d, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func TestNewOrdersHandler(t *testing.T) {
	o := NewOrdersHandler(&mocks.MockOrdersDataService{})

	assert.IsType(t, &OrdersHandler{}, o)
	assert.IsType(t, &mocks.MockOrdersDataService{}, o.OrdersDataSvc)
}

func TestCreateOrderSuccess(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	order, _ := json.Marshal(models.Order{
		Products: []models.Product{{
			Name:  "test-prod",
			Price: 100,
		}},
	})
	body := bytes.NewReader(order)
	c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/orders", body)
	mocks.CreateFunc = func(ctx context.Context, order interface{}) (*mongo.InsertOneResult, error) {
		data, err := ioutil.ReadFile("../../mockdata/createOrder.json")
		if err != nil {
			return nil, err
		}
		d, _ := UnMarshalCreateOrderResponse(data)
		return d, nil
	}

	// Call actual function
	o := NewOrdersHandler(&mocks.MockOrdersDataService{})
	o.Create(c)

	// Check results
	resp := w.Result()
	respBody, _ := io.ReadAll(resp.Body)
	respOrder, _ := UnMarshalCreateOrderResponse(respBody)
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
	assert.EqualValues(t, respOrder.InsertedID, "629fd50cb1e95cbe7ac12aae")
}

func TestCreateOrderFailure_DBError(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	order, _ := json.Marshal(models.Order{
		Products: []models.Product{{
			Name:  "test-prod",
			Price: 100,
		}},
	})
	body := bytes.NewReader(order)
	c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/orders", body)
	mocks.CreateFunc = func(ctx context.Context, order interface{}) (*mongo.InsertOneResult, error) {
		return nil, errors.New("db error")
	}

	// Call actual function
	o := NewOrdersHandler(&mocks.MockOrdersDataService{})
	o.Create(c)

	// Check results
	resp := w.Result()
	assert.EqualValues(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestCreateOrderFailure_BadRequest(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	order, _ := json.Marshal("Bad Request")
	body := bytes.NewReader(order)
	c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/orders", body)
	mocks.CreateFunc = func(ctx context.Context, order interface{}) (*mongo.InsertOneResult, error) {
		return nil, nil
	}

	// Call actual function
	o := NewOrdersHandler(&mocks.MockOrdersDataService{})
	o.Create(c)

	// Check results
	resp := w.Result()
	assert.EqualValues(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateOrderSuccess(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	id, _ := primitive.ObjectIDFromHex("629fd50cb1e95cbe7ac12aae")
	order, _ := json.Marshal(models.Order{
		ID: id,
		Products: []models.Product{{
			Name:  "test-prod",
			Price: 100,
		}},
	})
	body := bytes.NewReader(order)
	c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/orders", body)
	mocks.UpdateFunc = func(ctx context.Context, order interface{}) (int64, error) {
		return 1, nil
	}

	// Call actual function
	o := NewOrdersHandler(&mocks.MockOrdersDataService{})
	o.Create(c)

	// Check results
	resp := w.Result()
	respBody, _ := io.ReadAll(resp.Body)
	result, _ := strconv.Atoi(string(respBody))
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
	assert.EqualValues(t, 1, result)
}

func TestGetAllOrdersSuccess(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	mocks.GetAllFunc = func(ctx context.Context) (interface{}, error) {
		data, err := os.ReadFile("../../mockdata/allOrders.json")
		if err != nil {
			return nil, err
		}
		d, _ := UnMarshalOrdersResponse(data)
		return d, nil
	}

	// Call actual function
	o := NewOrdersHandler(&mocks.MockOrdersDataService{})
	o.GetAll(c)

	// Check results
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	orders, _ := UnMarshalOrdersResponse(body)
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
	assert.EqualValues(t, len(*orders), 100)
}

func TestGetAllOrdersFailure_DBRead(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	mocks.GetAllFunc = func(ctx context.Context) (interface{}, error) {
		_, err := os.ReadFile("../../mockdata/non-existing.json")
		return nil, err
	}

	// Call actual function
	o := NewOrdersHandler(&mocks.MockOrdersDataService{})
	o.GetAll(c)

	// Check results
	resp := w.Result()
	assert.EqualValues(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestGetOrderSuccess(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	const id = "629536b3fac02728de50c042"
	c.Params = []gin.Param{{Key: "id", Value: id}}
	mocks.GetByIdFunc = func(ctx context.Context, id string) (interface{}, error) {
		data, err := os.ReadFile("../../mockdata/order.json")
		if err != nil {
			return nil, err
		}
		d, _ := UnMarshalOrderResponse(data)
		return d, nil
	}

	// Call actual function
	o := NewOrdersHandler(&mocks.MockOrdersDataService{})
	o.GetByID(c)

	// Check results
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	order, _ := UnMarshalOrderResponse(body)
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
	assert.EqualValues(t, id, order.ID.Hex())
}

func TestGetOrderFailure_InvalidId(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	const id = ""
	c.Params = []gin.Param{{Key: "id", Value: id}}
	mocks.GetByIdFunc = func(ctx context.Context, id string) (interface{}, error) {
		data, err := os.ReadFile("../../mockdata/order.json")
		if err != nil {
			return nil, err
		}
		d, _ := UnMarshalOrderResponse(data)
		return d, nil
	}

	// Call actual function
	o := NewOrdersHandler(&mocks.MockOrdersDataService{})
	o.GetByID(c)

	// Check results
	resp := w.Result()
	assert.EqualValues(t, http.StatusBadRequest, resp.StatusCode)
}

func TestGetOrderFailure_DBRead(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	const id = "629536b3fac02728de50c042"
	c.Params = []gin.Param{{Key: "id", Value: id}}
	mocks.GetByIdFunc = func(ctx context.Context, id string) (interface{}, error) {
		_, err := os.ReadFile("../../mockdata/nan.json")
		return nil, err
	}

	// Call actual function
	o := NewOrdersHandler(&mocks.MockOrdersDataService{})
	o.GetByID(c)

	// Check results
	resp := w.Result()
	assert.EqualValues(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestDeleteOrderSuccess(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	const id = "629536b3fac02728de50c042"
	c.Params = []gin.Param{{Key: "id", Value: id}}
	mocks.DeleteByIdFunc = func(ctx context.Context, id string) (int64, error) {
		return 1, nil
	}

	// Call actual function
	o := NewOrdersHandler(&mocks.MockOrdersDataService{})
	o.DeleteByID(c)

	// Check results
	resp := w.Result()
	respBody, _ := io.ReadAll(resp.Body)
	result, _ := strconv.Atoi(string(respBody))
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
	assert.EqualValues(t, result, 1)
}

func TestDeleteOrderFailure_DBError(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	const id = "629536b3fac02728de50c042"
	c.Params = []gin.Param{{Key: "id", Value: id}}
	mocks.DeleteByIdFunc = func(ctx context.Context, id string) (int64, error) {
		return 1, errors.New("db error")
	}

	// Call actual function
	o := NewOrdersHandler(&mocks.MockOrdersDataService{})
	o.DeleteByID(c)

	// Check results
	resp := w.Result()
	assert.EqualValues(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestDeleteOrderFailure_BadRequest(t *testing.T) {
	// Test Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	const id = ""
	c.Params = []gin.Param{{Key: "id", Value: id}}
	mocks.DeleteByIdFunc = func(ctx context.Context, id string) (int64, error) {
		return 0, nil
	}

	// Call actual function
	o := NewOrdersHandler(&mocks.MockOrdersDataService{})
	o.DeleteByID(c)

	// Check results
	resp := w.Result()
	assert.EqualValues(t, http.StatusBadRequest, resp.StatusCode)
}
*/
