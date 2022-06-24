package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/mocks"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
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
	o := NewOrdersController(&mocks.MockDataService{})

	assert.IsType(t, &OrdersController{}, o)
	assert.IsType(t, &mocks.MockDataService{}, o.dataSvc)
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
	c.Request, _ = http.NewRequest("POST", "/api/v1/orders", body)
	mocks.CreateFunc = func(interface{}) (*mongo.InsertOneResult, error) {
		data, err := ioutil.ReadFile("../mockdata/createOrder.json")
		if err != nil {
			return nil, err
		}
		d, _ := UnMarshalCreateOrderResponse(data)
		return d, nil
	}

	// Call actual function
	o := NewOrdersController(&mocks.MockDataService{})
	o.Post(c)

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
	c.Request, _ = http.NewRequest("POST", "/api/v1/orders", body)
	mocks.CreateFunc = func(interface{}) (*mongo.InsertOneResult, error) {
		return nil, errors.New("db error")
	}

	// Call actual function
	o := NewOrdersController(&mocks.MockDataService{})
	o.Post(c)

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
	c.Request, _ = http.NewRequest("POST", "/api/v1/orders", body)
	mocks.CreateFunc = func(interface{}) (*mongo.InsertOneResult, error) {
		return nil, nil
	}

	// Call actual function
	o := NewOrdersController(&mocks.MockDataService{})
	o.Post(c)

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
	c.Request, _ = http.NewRequest("POST", "/api/v1/orders", body)
	mocks.UpdateFunc = func(interface{}) (int64, error) {
		return 1, nil
	}

	// Call actual function
	o := NewOrdersController(&mocks.MockDataService{})
	o.Post(c)

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
	mocks.GetAllFunc = func() (interface{}, error) {
		data, err := ioutil.ReadFile("../mockdata/allOrders.json")
		if err != nil {
			return nil, err
		}
		d, _ := UnMarshalOrdersResponse(data)
		return d, nil
	}

	// Call actual function
	o := NewOrdersController(&mocks.MockDataService{})
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
	mocks.GetAllFunc = func() (interface{}, error) {
		_, err := ioutil.ReadFile("../mockdata/non-existing.json")
		return nil, err
	}

	// Call actual function
	o := NewOrdersController(&mocks.MockDataService{})
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
	mocks.GetByIdFunc = func(id string) (interface{}, error) {
		data, err := ioutil.ReadFile("../mockdata/order.json")
		if err != nil {
			return nil, err
		}
		d, _ := UnMarshalOrderResponse(data)
		return d, nil
	}

	// Call actual function
	o := NewOrdersController(&mocks.MockDataService{})
	o.GetById(c)

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
	mocks.GetByIdFunc = func(id string) (interface{}, error) {
		data, err := ioutil.ReadFile("../mockdata/order.json")
		if err != nil {
			return nil, err
		}
		d, _ := UnMarshalOrderResponse(data)
		return d, nil
	}

	// Call actual function
	o := NewOrdersController(&mocks.MockDataService{})
	o.GetById(c)

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
	mocks.GetByIdFunc = func(id string) (interface{}, error) {
		_, err := ioutil.ReadFile("../mockdata/nan.json")
		return nil, err
	}

	// Call actual function
	o := NewOrdersController(&mocks.MockDataService{})
	o.GetById(c)

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
	mocks.DeleteByIdFunc = func(id string) (int64, error) {
		return 1, nil
	}

	// Call actual function
	o := NewOrdersController(&mocks.MockDataService{})
	o.DeleteById(c)

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
	mocks.DeleteByIdFunc = func(id string) (int64, error) {
		return 1, errors.New("db error")
	}

	// Call actual function
	o := NewOrdersController(&mocks.MockDataService{})
	o.DeleteById(c)

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
	mocks.DeleteByIdFunc = func(id string) (int64, error) {
		return 0, nil
	}

	// Call actual function
	o := NewOrdersController(&mocks.MockDataService{})
	o.DeleteById(c)

	// Check results
	resp := w.Result()
	assert.EqualValues(t, http.StatusBadRequest, resp.StatusCode)
}
