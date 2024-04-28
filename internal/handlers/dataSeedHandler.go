package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/data"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
)

const (
	seedRecordCount = 10000
)

type SeedHandler struct {
	oDataSvc db.OrdersDataService
}

func NewDataSeedHandler(svc db.OrdersDataService) *SeedHandler {
	sc := &SeedHandler{
		oDataSvc: svc,
	}
	return sc
}

func (s *SeedHandler) SeedDB(c *gin.Context) {
	for i := 0; i < seedRecordCount; i++ {
		products := []data.Product{
			{
				Name:      faker.Name(),
				Price:     util.RandomPrice(),
				UpdatedAt: time.Now(),
			},
			{
				Name:      faker.Name(),
				Price:     util.RandomPrice(),
				UpdatedAt: time.Now(),
			},
		}

		po := &data.Order{
			Version:     1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Products:    products,
			User:        faker.Email(),
			Status:      data.OrderPending,
			TotalAmount: util.CalculateTotalAmount(products),
		}

		_, err := s.oDataSvc.Create(c, po)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "failed to insert data",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully inserted fake data",
		"Count":   seedRecordCount,
	})
}
