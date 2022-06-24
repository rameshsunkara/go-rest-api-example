package mocks

import "go.mongodb.org/mongo-driver/mongo"

var (
	PingFunc       func() error
	CreateFunc     func(purchaseOrder interface{}) (*mongo.InsertOneResult, error)
	UpdateFunc     func(purchaseOrder interface{}) (int64, error)
	GetAllFunc     func() (interface{}, error)
	GetByIdFunc    func(id string) (interface{}, error)
	DeleteByIdFunc func(id string) (int64, error)
)

type MockDataMgr struct{}

func (m *MockDataMgr) Ping() error {
	return PingFunc()
}

func (m *MockDataMgr) Client() (*mongo.Client, error) {
	return nil, nil
}

func (m *MockDataMgr) Database() (*mongo.Database, error) {
	return nil, nil
}

type MockDataService struct{}

func (m *MockDataService) Create(purchaseOrder interface{}) (*mongo.InsertOneResult, error) {
	return CreateFunc(purchaseOrder)
}

func (m *MockDataService) Update(purchaseOrder interface{}) (int64, error) {
	return UpdateFunc(purchaseOrder)
}

func (m *MockDataService) GetAll() (interface{}, error) {
	return GetAllFunc()
}

func (m *MockDataService) GetById(id string) (interface{}, error) {
	return GetByIdFunc(id)
}

func (m *MockDataService) DeleteById(id string) (int64, error) {
	return DeleteByIdFunc(id)
}
