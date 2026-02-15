package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/bogdanutanu/go-rest-api-example/internal/db"
	"github.com/bogdanutanu/go-rest-api-example/internal/models/data"
	"github.com/bogdanutanu/go-rest-api-example/internal/utilities"
	"github.com/bogdanutanu/go-rest-api-example/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
)

const (
	seedRecordCount = 10000
)

type SeedHandler struct {
	oDataSvc db.OrdersDataService
	lgr      logger.Logger
}

func NewDataSeedHandler(lgr logger.Logger, svc db.OrdersDataService) (*SeedHandler, error) {
	if lgr == nil || svc == nil {
		return nil, errors.New("failed to create local DB seed handler")
	}
	return &SeedHandler{
		oDataSvc: svc,
		lgr:      lgr,
	}, nil
}

func (s *SeedHandler) SeedDB(c *gin.Context) {
	for i := 0; i < seedRecordCount; i++ {
		products := []data.Product{
			{
				Name:      faker.Name(),
				Price:     utilities.RandomPrice(),
				UpdatedAt: time.Now(),
			},
			{
				Name:      faker.Name(),
				Price:     utilities.RandomPrice(),
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
			TotalAmount: utilities.CalculateTotalAmount(products),
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
