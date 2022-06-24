package mocks

import (
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func (m *MockDataMgr) Database() (db.MongoDatabase, error) {
	return &MockMongoDataBase{}, nil
}

type MockMongoDataBase struct{}

func (m *MockMongoDataBase) Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return nil
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
