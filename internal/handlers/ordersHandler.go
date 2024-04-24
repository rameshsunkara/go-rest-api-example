package handlers

import (
	"net/http"
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
	orders, err := o.oDataSvc.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"message": "error occurred while retrieving purchase orders", "error": err})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, orders)
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
