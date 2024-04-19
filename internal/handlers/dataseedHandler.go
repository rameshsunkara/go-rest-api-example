package handlers

import (
	"math/rand"
	"net/http"

	"github.com/go-faker/faker/v4"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/types"

	"github.com/gin-gonic/gin"
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
		product := []types.Product{
			{
				Name:      faker.Name(),
				Price:     (uint)(rand.Intn(90) + 10),
				Remarks:   faker.Sentence(),
				UpdatedAt: faker.TimeString(),
			},
			{
				Name:      faker.Name(),
				Price:     (uint)(rand.Intn(1000) + 10),
				Remarks:   faker.Sentence(),
				UpdatedAt: faker.TimeString(),
			},
		}

		po := &types.Order{
			Products: product,
		}
		_, err := s.OrdersDataSvc.Create(c, po)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable inserted data",
			})
			panic("Unable to insert data")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully inserted fake data",
		"Count":   seedRecordCount,
	})
}
