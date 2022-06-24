package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
)

const (
	OrderIdPath = "id" // Request path variable
)

type OrdersController struct {
	dataSvc db.DataService
}

func NewOrdersController(svc db.DataService) *OrdersController {
	ic := &OrdersController{
		dataSvc: svc,
	}
	return ic
}

// Post  godoc
// @Summary      Creates or Updates an order
// @Description  Used to either create or update an order
// @Tags         Fetch
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /orders/ [post]
func (oHandler *OrdersController) Post(c *gin.Context) {
	purchaseRequest := models.Order{}

	if err := c.BindJSON(&purchaseRequest); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if purchaseRequest.ID.IsZero() {
		if uid, _ := oHandler.dataSvc.Create(&purchaseRequest); uid != nil {
			c.JSON(http.StatusOK, uid)
			return
		}
	} else {
		if updatedCount, _ := oHandler.dataSvc.Update(&purchaseRequest); updatedCount != 0 {
			c.JSON(http.StatusOK, updatedCount)
			return
		}
	}

	c.JSON(http.StatusInternalServerError, "Unexpected Error occurred")
}

// GetAll  godoc
// @Summary      Fetch all orders
// @Description  Fetches all orders
// @Tags         Fetch
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /orders/ [get]
func (oHandler *OrdersController) GetAll(c *gin.Context) {
	orders, err := oHandler.dataSvc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error occurred while retrieved purchase orders", "error": err})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, orders)
}

// GetById  godoc
// @Summary      Fetch single Order document identified by give id
// @Description  Fetch single Order document identified by give id
// @Param        id   path      string  true  "Order ID"
// @Tags         Fetch
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      500            {string}  string  "bad request"
// @Router       /orders/{id} [get]
func (oHandler *OrdersController) GetById(c *gin.Context) {
	id := c.Param(OrderIdPath)
	if id != "" {
		order, err := oHandler.dataSvc.GetById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error to retrieve order details", "error": err.Error()})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, order)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
	c.Abort()
}

// DeleteById  godoc
// @Summary      Delete single Order document identified by give id
// @Description  Delete single Order document identified by give id
// @Param        id   path      string  true  "Order ID"
// @Tags         Fetch
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      500            {string}  string  "bad request"
// @Router       /orders/{id} [delete]
func (oHandler *OrdersController) DeleteById(c *gin.Context) {
	id := c.Param(OrderIdPath)
	if id != "" {
		count, err := oHandler.dataSvc.DeleteById(id)
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
