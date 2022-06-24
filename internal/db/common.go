package db

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

type DataManager interface {
	Database() (MongoDatabase, error)
	Ping() error
}

type DataService interface {
	Create(purchaseOrder interface{}) (*mongo.InsertOneResult, error)
	Update(purchaseOrder interface{}) (int64, error)
	GetAll() (interface{}, error)
	GetById(id string) (interface{}, error)
	DeleteById(id string) (int64, error)
}
