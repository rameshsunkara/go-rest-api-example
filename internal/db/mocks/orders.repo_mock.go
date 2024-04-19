package mocks

import (
	"context"

	"github.com/rameshsunkara/go-rest-api-example/internal/types"
)

var (
	CreateFunc     func(ctx context.Context, purchaseOrder *types.Order) (string, error)
	UpdateFunc     func(ctx context.Context, purchaseOrder *types.Order) error
	GetAllFunc     func(ctx context.Context) (*[]types.Order, error)
	GetByIDFunc    func(ctx context.Context, id string) (*types.Order, error)
	DeleteByIDFunc func(ctx context.Context, id string) (int64, error)
)

type MockOrdersDataService struct{}

func (m *MockOrdersDataService) Create(ctx context.Context, purchaseOrder *types.Order) (string, error) {
	return CreateFunc(ctx, purchaseOrder)
}

func (m *MockOrdersDataService) Update(ctx context.Context, purchaseOrder *types.Order) error {
	return UpdateFunc(ctx, purchaseOrder)
}

func (m *MockOrdersDataService) GetAll(ctx context.Context) (*[]types.Order, error) {
	return GetAllFunc(ctx)
}

func (m *MockOrdersDataService) GetByID(ctx context.Context, id string) (*types.Order, error) {
	return GetByIDFunc(ctx, id)
}

func (m *MockOrdersDataService) DeleteByID(ctx context.Context, id string) (int64, error) {
	return DeleteByIDFunc(ctx, id)
}
