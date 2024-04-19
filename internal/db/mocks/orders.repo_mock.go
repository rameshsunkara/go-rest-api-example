package mocks

import (
	"context"

	"github.com/rameshsunkara/go-rest-api-example/internal/types"
)

type MockOrdersDataService struct {
	CreateFunc     func(ctx context.Context, purchaseOrder *types.Order) (string, error)
	UpdateFunc     func(ctx context.Context, purchaseOrder *types.Order) error
	GetAllFunc     func(ctx context.Context) (*[]types.Order, error)
	GetByIDFunc    func(ctx context.Context, id string) (*types.Order, error)
	DeleteByIDFunc func(ctx context.Context, id string) (int64, error)
}

func (m *MockOrdersDataService) Create(ctx context.Context, purchaseOrder *types.Order) (string, error) {
	return m.CreateFunc(ctx, purchaseOrder)
}

func (m *MockOrdersDataService) Update(ctx context.Context, purchaseOrder *types.Order) error {
	return m.UpdateFunc(ctx, purchaseOrder)
}

func (m *MockOrdersDataService) GetAll(ctx context.Context) (*[]types.Order, error) {
	return m.GetAllFunc(ctx)
}

func (m *MockOrdersDataService) GetByID(ctx context.Context, id string) (*types.Order, error) {
	return m.GetByIDFunc(ctx, id)
}

func (m *MockOrdersDataService) DeleteByID(ctx context.Context, id string) (int64, error) {
	return m.DeleteByIDFunc(ctx, id)
}
