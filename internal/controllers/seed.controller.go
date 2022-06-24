package controllers

import (
	"math/rand"
	"net/http"

	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
)

const (
	SeedRecordCount = 500
)

type SeedController struct {
	dataSvc db.DataService
}

func NewSeedController(svc db.DataService) *SeedController {
	ic := &SeedController{
		dataSvc: svc,
	}
	return ic
}

func (s *SeedController) SeedDB(c *gin.Context) {
	for i := 0; i < SeedRecordCount; i++ {
		product := []models.Product{
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

		po := &models.Order{
			Products: product,
		}
		_, err := s.dataSvc.Create(po)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable inserted data",
			})
			panic("Unable to insert data")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully inserted fake data",
		"Count":   SeedRecordCount,
	})
}
