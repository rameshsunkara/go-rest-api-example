package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/types"
)

const (
	OrderIdPath = "id" // Request path variable
)

type OrdersController struct {
	dataSvc db.OrdersDataService
}

func NewOrdersController(svc db.OrdersDataService) *OrdersController {
	ic := &OrdersController{
		dataSvc: svc,
	}
	return ic
}

func (oHandler *OrdersController) Post(c *gin.Context) {
	purchaseRequest := types.Order{}

	if err := c.ShouldBind(&purchaseRequest); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if purchaseRequest.ID.IsZero() {
		if uid, _ := oHandler.dataSvc.Create(c, &purchaseRequest); uid != nil {
			c.JSON(http.StatusOK, uid)
			return
		}
	} else {
		if updatedCount, _ := oHandler.dataSvc.Update(c, &purchaseRequest); updatedCount != 0 {
			c.JSON(http.StatusOK, updatedCount)
			return
		}
	}

	c.JSON(http.StatusInternalServerError, "Unexpected Error occurred")
}

func (oHandler *OrdersController) GetAll(c *gin.Context) {
	orders, err := oHandler.dataSvc.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error occurred while retrieved purchase orders", "error": err})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (oHandler *OrdersController) GetById(c *gin.Context) {
	id := c.Param(OrderIdPath)
	if id != "" {
		order, err := oHandler.dataSvc.GetById(c, id)
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

func (oHandler *OrdersController) DeleteById(c *gin.Context) {
	id := c.Param(OrderIdPath)
	if id != "" {
		count, err := oHandler.dataSvc.DeleteById(c, id)
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
