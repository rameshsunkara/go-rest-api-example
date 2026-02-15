package handlers

import (
	errors2 "errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bogdanutanu/go-rest-api-example/internal/db"
	"github.com/bogdanutanu/go-rest-api-example/internal/errors"
	"github.com/bogdanutanu/go-rest-api-example/internal/models/data"
	"github.com/bogdanutanu/go-rest-api-example/internal/models/external"
	"github.com/bogdanutanu/go-rest-api-example/internal/utilities"
	"github.com/bogdanutanu/go-rest-api-example/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	OrderIDPath = "id"
	MaxPageSize = 100
)

// OrdersHandler handles order-related HTTP requests.
type OrdersHandler struct {
	oDataSvc db.OrdersDataService
	logger   logger.Logger
}

// NewOrdersHandler creates a new OrdersHandler.
func NewOrdersHandler(lgr logger.Logger, dSvc db.OrdersDataService) (*OrdersHandler, error) {
	if lgr == nil || dSvc == nil {
		return nil, errors2.New("missing required parameters to create orders handler")
	}
	return &OrdersHandler{oDataSvc: dSvc, logger: lgr}, nil
}

// Create handles POST /orders.
func (o *OrdersHandler) Create(c *gin.Context) {
	lgr, requestID := o.logger.WithReqID(c)
	var orderInput external.OrderInput
	if err := c.ShouldBindJSON(&orderInput); err != nil {
		o.abortWithAPIError(c, lgr, http.StatusBadRequest, errors.OrderCreateInvalidInput,
			"Invalid order request body", requestID, err)
		return
	}

	products := make([]data.Product, len(orderInput.Products))
	for i, p := range orderInput.Products {
		products[i] = data.Product{Name: p.Name, Price: p.Price, Quantity: p.Quantity}
	}

	order := data.Order{
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Products:    products,
		User:        faker.Email(), // TODO: Replace with actual user email from trusted source such as JWT token
		TotalAmount: utilities.CalculateTotalAmount(products),
		Status:      data.OrderPending,
	}

	id, err := o.oDataSvc.Create(c, &order)
	if err != nil {
		o.abortWithAPIError(c, lgr, http.StatusInternalServerError, errors.OrderCreateServerError,
			errors.UnexpectedErrorMessage, requestID, err)
		return
	}

	extOrder := external.Order{
		ID:          id,
		CreatedAt:   utilities.FormatTimeToISO(order.CreatedAt),
		UpdatedAt:   utilities.FormatTimeToISO(order.UpdatedAt),
		Products:    order.Products,
		User:        order.User,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		Version:     order.Version,
	}
	c.JSON(http.StatusCreated, extOrder)
}

// GetAll handles GET /orders.
func (o *OrdersHandler) GetAll(c *gin.Context) {
	lgr, requestID := o.logger.WithReqID(c)
	limit, apiErr := o.parseLimitQueryParam(c)
	if apiErr != nil {
		c.AbortWithStatusJSON(apiErr.HTTPStatusCode, apiErr)
		return
	}

	orders, err := o.oDataSvc.GetAll(c, limit)
	if err != nil {
		o.abortWithAPIError(c, lgr, http.StatusInternalServerError, errors.OrdersGetServerError,
			errors.UnexpectedErrorMessage, requestID, err)
		return
	}

	extOrders := make([]external.Order, 0, len(*orders))
	for _, o := range *orders {
		extOrders = append(extOrders, external.Order{
			ID:          o.ID.Hex(),
			Version:     o.Version,
			Status:      o.Status,
			TotalAmount: o.TotalAmount,
			User:        o.User,
			CreatedAt:   utilities.FormatTimeToISO(o.CreatedAt),
			UpdatedAt:   utilities.FormatTimeToISO(o.UpdatedAt),
			Products:    o.Products,
		})
	}
	c.JSON(http.StatusOK, extOrders)
}

// GetByID handles GET /orders/:id.
func (o *OrdersHandler) GetByID(c *gin.Context) {
	lgr, requestID := o.logger.WithReqID(c)
	id := c.Param(OrderIDPath)
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil || oID.IsZero() {
		o.abortWithAPIError(c, lgr, http.StatusBadRequest, errors.OrderGetInvalidParams, "invalid order ID", requestID, err)
		return
	}
	order, err := o.oDataSvc.GetByID(c, oID)
	if err != nil {
		if errors2.Is(err, db.ErrPOIDNotFound) {
			o.abortWithAPIError(c, lgr, http.StatusNotFound, errors.OrderGetNotFound,
				"order not found", requestID, err)
			return
		}
		o.abortWithAPIError(c, lgr, http.StatusInternalServerError, errors.OrdersGetServerError,
			"failed to fetch order", requestID, err)
		return
	}
	c.JSON(http.StatusOK, order)
}

// DeleteByID handles DELETE /orders/:id.
func (o *OrdersHandler) DeleteByID(c *gin.Context) {
	lgr, requestID := o.logger.WithReqID(c)
	id := c.Param(OrderIDPath)
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil || oID.IsZero() {
		o.abortWithAPIError(c, lgr, http.StatusBadRequest, errors.OrderDeleteInvalidID, "invalid order ID", requestID, err)
		return
	}
	if dbErr := o.oDataSvc.DeleteByID(c, oID); dbErr != nil {
		if errors2.Is(dbErr, db.ErrPOIDNotFound) {
			o.abortWithAPIError(c, lgr, http.StatusNotFound, errors.OrderDeleteNotFound,
				"could not delete order", requestID, dbErr)
			return
		}
		o.abortWithAPIError(c, lgr, http.StatusInternalServerError, errors.OrderDeleteServerError,
			"could not delete order", requestID, dbErr)
		return
	}
	c.Status(http.StatusNoContent)
}

// parseLimitQueryParam parses and validates the "limit" query parameter.
func (o *OrdersHandler) parseLimitQueryParam(c *gin.Context) (int64, *external.APIError) {
	lgr, requestID := o.logger.WithReqID(c)
	l := db.DefaultPageSize
	if input, exists := c.GetQuery("limit"); exists && input != "" {
		val, err := strconv.Atoi(input)
		if err != nil || val < 1 || val > MaxPageSize {
			apiErr := &external.APIError{
				HTTPStatusCode: http.StatusBadRequest,
				ErrorCode:      "",
				Message:        fmt.Sprintf("Integer value within 1 and %d is expected for limit query param", MaxPageSize),
				DebugID:        requestID,
			}
			lgr.Error().
				Int("HttpStatusCode", apiErr.HTTPStatusCode).
				Str("ErrorCode", apiErr.ErrorCode).
				Msg(apiErr.Message)
			return 0, apiErr
		}
		l = val
	}
	return int64(l), nil
}

// abortWithAPIError logs and aborts the request with a standardized API error response.
func (o *OrdersHandler) abortWithAPIError(
	c *gin.Context,
	lgr logger.Logger,
	status int,
	errorCode, message, debugID string,
	err error,
) {
	apiErr := &external.APIError{
		HTTPStatusCode: status,
		ErrorCode:      errorCode,
		Message:        message,
		DebugID:        debugID,
	}
	event := lgr.Error().Int("HttpStatusCode", status).Str("ErrorCode", errorCode)
	if err != nil {
		event = event.Err(err)
	}
	event.Msg(message)
	c.AbortWithStatusJSON(status, apiErr)
}
