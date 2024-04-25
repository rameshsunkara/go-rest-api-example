package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/errors"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/data"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/external"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
)

const (
	OrderIDPath = "id" // Request path variable
	MaxPageSize = 100  // Maximum number of records that can be fetched in a single request
)

type OrdersHandler struct {
	oDataSvc db.OrdersDataService
	logger   *logger.AppLogger
}

func NewOrdersHandler(dSvc db.OrdersDataService, lgr *logger.AppLogger) *OrdersHandler {
	ic := &OrdersHandler{
		oDataSvc: dSvc,
		logger:   lgr,
	}
	return ic
}

func (o *OrdersHandler) Create(c *gin.Context) {
	lgr, requestID := o.logger.WithReqID(c)
	// Validate  inputs : fail fast order
	// Parse request body
	var orderInput external.OrderInput
	if err := c.ShouldBindJSON(&orderInput); err != nil {
		apiErr := &external.APIError{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorCode:      errors.OrderCreateInvalidInput,
			Message:        "Invalid order request body",
			DebugID:        requestID,
		}
		lgr.Error().
			Err(err).
			Int("HttpStatusCode", apiErr.HTTPStatusCode).
			Str("ErrorCode", apiErr.ErrorCode).
			Msg(apiErr.Message)
		c.JSON(apiErr.HTTPStatusCode, apiErr)
		return
	}

	// Convert ProductInput to Product
	var products []data.Product
	for _, productInput := range orderInput.Products {
		product := data.Product{
			Name:     productInput.Name,
			Price:    productInput.Price,
			Quantity: productInput.Quantity,
		}
		products = append(products, product)
	}

	// Create new Order object
	order := data.Order{
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Products:    products,
		User:        faker.Email(), // TODO: Replace with actual user email from trusted source such as JWT token
		TotalAmount: util.CalculateTotalAmount(products),
		Status:      data.OrderPending,
	}

	if id, err := o.oDataSvc.Create(c, &order); err == nil {
		// Return success response
		extOrder := external.Order{
			ID:          id,
			CreatedAt:   util.FormatTimeToISO(order.CreatedAt),
			UpdatedAt:   util.FormatTimeToISO(order.UpdatedAt),
			Products:    order.Products,
			User:        order.User,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			Version:     order.Version,
		}
		c.JSON(http.StatusCreated, extOrder)
		return
	}

	apiErr := &external.APIError{
		HTTPStatusCode: http.StatusInternalServerError,
		ErrorCode:      errors.OrderCreateServerError,
		Message:        errors.UnexpectedErrorMessage,
		DebugID:        requestID,
	}
	lgr.Error().
		Int("HttpStatusCode", apiErr.HTTPStatusCode).
		Str("ErrorCode", apiErr.ErrorCode).
		Msg(apiErr.Message)
	c.JSON(apiErr.HTTPStatusCode, apiErr)
}

func (o *OrdersHandler) GetAll(c *gin.Context) {
	lgr, requestID := o.logger.WithReqID(c)
	// Validate  inputs : fail fast order
	orders, err := o.oDataSvc.GetAll(c)
	var extOrders []external.Order
	if orders != nil {
		extOrders = make([]external.Order, len(*orders))
		for i, o := range *orders {
			extOrders[i] = external.Order{
				ID:          o.ID.Hex(),
				Version:     o.Version,
				Status:      o.Status,
				TotalAmount: o.TotalAmount,
				User:        o.User,
				CreatedAt:   util.FormatTimeToISO(o.CreatedAt),
				UpdatedAt:   util.FormatTimeToISO(o.UpdatedAt),
				Products: o.Products,
			}
		}
	}

	if err != nil {
		apiErr := &external.APIError{
			HTTPStatusCode: http.StatusInternalServerError,
			ErrorCode:      errors.OrdersGetServerError,
			Message:        errors.UnexpectedErrorMessage,
			DebugID:        requestID,
		}
		lgr.Error().
			Int("HttpStatusCode", apiErr.HTTPStatusCode).
			Str("ErrorCode", apiErr.ErrorCode).
			Msg(apiErr.Message)
		c.JSON(apiErr.HTTPStatusCode, apiErr)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, extOrders)
}

func (o *OrdersHandler) GetByID(c *gin.Context) {
	id := c.Param(OrderIDPath)
	if id != "" {
		order, err := o.oDataSvc.GetByID(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"message": "Error to retrieve order details", "error": err.Error()})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, order)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
	c.Abort()
}

func (o *OrdersHandler) DeleteByID(c *gin.Context) {
	id := c.Param(OrderIDPath)
	if id != "" {
		count, err := o.oDataSvc.DeleteByID(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error to retrieve order details", "error": err.Error()})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, count)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
	c.Abort()
}

func (o *OrdersHandler) parseLimitQueryParam(c *gin.Context, selectParam bool) (int, *external.APIError) {
	lgr, requestId := o.logger.WithReqID(c)
	l := db.DefaultPageSize
	if input, exists := c.GetQuery("limit"); exists && input != "" {
		if selectParam {
			apiErr := &external.APIError{
				HTTPStatusCode: http.StatusBadRequest,
				ErrorCode:      "",
				Message:        "limit query param is not supported when select query param is used",
				DebugID:        requestId,
			}
			lgr.Error().
				Int("HttpStatusCode", apiErr.HTTPStatusCode).
				Str("ErrorCode", apiErr.ErrorCode).
				Msg(apiErr.Message)
			return 0, apiErr
		}
		var err error
		l, err = strconv.Atoi(input)
		if err != nil {
			apiErr := &external.APIError{
				HTTPStatusCode: http.StatusBadRequest,
				ErrorCode:      "",
				Message:        fmt.Sprintf("Integer value within 1 and %d is expected for limit query param",
					MaxPageSize),
				DebugID:        requestId,
			}
			lgr.Error().
				Int("HttpStatusCode", apiErr.HTTPStatusCode).
				Str("ErrorCode", apiErr.ErrorCode).
				Msg(apiErr.Message)
			return 0, apiErr
		}
		if l < 1 || l > MaxPageSize {
			apiErr := &external.APIError{
				HTTPStatusCode: http.StatusBadRequest,
				ErrorCode:      "",
				Message:        fmt.Sprintf("Integer value within 1 and %d is expected for limit query param",
					MaxPageSize),
				DebugID:        requestId,
			}
			lgr.Error().
				Int("HttpStatusCode", apiErr.HTTPStatusCode).
				Str("ErrorCode", apiErr.ErrorCode).
				Msg(apiErr.Message)
			return 0, apiErr
		}
	}
	return l, nil
}