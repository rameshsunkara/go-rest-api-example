package mocks

import (
	"context"

	"github.com/bogdanutanu/go-rest-api-example/internal/models/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockOrdersDataService struct {
	CreateFunc func(ctx context.Context, purchaseOrder *data.Order) (string, error)
	// UpdateFunc     func(ctx context.Context, purchaseOrder *data.Order) error
	GetAllFunc     func(ctx context.Context, limit int64) (*[]data.Order, error)
	GetByIDFunc    func(ctx context.Context, id primitive.ObjectID) (*data.Order, error)
	DeleteByIDFunc func(ctx context.Context, id primitive.ObjectID) error
}

func (m *MockOrdersDataService) Create(ctx context.Context, purchaseOrder *data.Order) (string, error) {
	return m.CreateFunc(ctx, purchaseOrder)
}

func (m *MockOrdersDataService) Update(_ context.Context, _ *data.Order) error {
	return nil
}

func (m *MockOrdersDataService) GetAll(ctx context.Context, limit int64) (*[]data.Order, error) {
	return m.GetAllFunc(ctx, limit)
}

func (m *MockOrdersDataService) GetByID(ctx context.Context, id primitive.ObjectID) (*data.Order, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m *MockOrdersDataService) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	return m.DeleteByIDFunc(ctx, id)
}
