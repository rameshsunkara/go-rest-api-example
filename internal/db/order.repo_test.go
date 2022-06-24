package db

import (
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllSuccess(t *testing.T) {
	d, _ := dbMgr.Database()
	dSvc := NewOrderDataService(d)
	results, _ := dSvc.GetAll()
	orders := results.(*[]models.Order)
	assert.EqualValues(t, 100, len(*orders))
}

func TestGetByIdSuccess_NoData(t *testing.T) {
	d, _ := dbMgr.Database()
	dSvc := NewOrderDataService(d)
	const id = "hola-non-id"
	result, _ := dSvc.GetById(id)
	assert.Nil(t, result)
}
