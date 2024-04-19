package handlers

import (
	"crypto/rand"
	"math/big"
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
	var price *big.Int
	var err error
	for i := 0; i < seedRecordCount; i++ {
		if price, err = rand.Int(rand.Reader, big.NewInt(1000)); err != nil {
			// default price
			price = big.NewInt(100)
		}
		products := []types.Product{
			{
				Name:        faker.Name(),
				Price:       price.Uint64(),
				Description: faker.Sentence(),
				UpdatedAt:   faker.TimeString(),
			},
			{
				Name:        faker.Name(),
				Price:       price.Uint64(),
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
			panic("failed to insert data")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully inserted fake data",
		"Count":   seedRecordCount,
	})
}
