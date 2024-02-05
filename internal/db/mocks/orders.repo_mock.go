package mocks

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	CreateFunc     func(ctx context.Context, purchaseOrder interface{}) (*mongo.InsertOneResult, error)
	UpdateFunc     func(ctx context.Context, purchaseOrder interface{}) (int64, error)
	GetAllFunc     func(ctx context.Context) (interface{}, error)
	GetByIdFunc    func(ctx context.Context, id string) (interface{}, error)
	DeleteByIdFunc func(ctx context.Context, id string) (int64, error)
)

type MockOrdersDataService struct{}

func (m *MockOrdersDataService) Create(ctx context.Context, purchaseOrder interface{}) (*mongo.InsertOneResult, error) {
	return CreateFunc(ctx, purchaseOrder)
}

func (m *MockOrdersDataService) Update(ctx context.Context, purchaseOrder interface{}) (int64, error) {
	return UpdateFunc(ctx, purchaseOrder)
}

func (m *MockOrdersDataService) GetAll(ctx context.Context) (interface{}, error) {
	return GetAllFunc(ctx)
}

func (m *MockOrdersDataService) GetById(ctx context.Context, id string) (interface{}, error) {
	return GetByIdFunc(ctx, id)
}

func (m *MockOrdersDataService) DeleteById(ctx context.Context, id string) (int64, error) {
	return DeleteByIdFunc(ctx, id)
}
