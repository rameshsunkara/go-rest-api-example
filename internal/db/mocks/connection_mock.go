package mocks

import (
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockMongoMgr struct {
	PingFunc func() error
}

func (m *MockMongoMgr) Ping() error {
	return m.PingFunc()
}

func (m *MockMongoMgr) Database() db.MongoDatabase {
	return &MockMongoDataBase{}
}

func (m *MockMongoMgr) Disconnect() error {
	return nil
}

type MockMongoDataBase struct{}

func (m *MockMongoDataBase) Collection(_ string, _ ...*options.CollectionOptions) *mongo.Collection {
	return nil
}
