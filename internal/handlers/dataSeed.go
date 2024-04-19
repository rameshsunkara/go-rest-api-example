package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/types"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
)

const (
	seedRecordCount = 10000
)

type SeedController struct {
	OrdersDataSvc db.OrdersDataService
}

func NewSeedController(svc db.OrdersDataService) *SeedController {
	ic := &SeedController{
		OrdersDataSvc: svc,
	}
	return ic
}

func (s *SeedController) SeedDB(c *gin.Context) {
	for i := 0; i < seedRecordCount; i++ {
		products := []types.Product{
			{
				Name:        faker.Name(),
				Price:       util.RandomPrice(),
				Description: faker.Sentence(),
				UpdatedAt:   faker.TimeString(),
			},
			{
				Name:        faker.Name(),
				Price:       util.RandomPrice(),
				Description: faker.Sentence(),
				UpdatedAt:   faker.TimeString(),
			},
		}

		po := &types.Order{
			CreatedAt: util.CurrentISOTime(),
			UpdatedAt: util.CurrentISOTime(),
			Products:  products,
			User:      faker.Email(),
		}

		_, err := s.OrdersDataSvc.Create(c, po)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
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
