// Package mongodb provides MongoDB connection utilities and abstractions.
// It offers a production-ready, idiomatic Go API for MongoDB operations with
// features like connection pooling, credential management, and flexible configuration.
package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabase defines the interface for MongoDB database operations.
type MongoDatabase interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

// MongoManager defines the interface for MongoDB connection management.
type MongoManager interface {
	Database() MongoDatabase
	DatabaseByName(name string) MongoDatabase
	Ping() error
	Disconnect() error
}

// MongoCredentials represents MongoDB authentication credentials.
type MongoCredentials struct {
	Username string `json:"username,omitempty" log:"-"`
	Password string `json:"password,omitempty" log:"-"`
}

// MongoOptions represents optional MongoDB connection settings
// Note: Hosts should be provided in "hostname:port" format directly to the ConnectionURL function.
type MongoOptions struct {
	UseSRV         bool   `json:"useSRV,omitempty"`         // Use SRV connection
	ReplicaSet     string `json:"replicaSet,omitempty"`     // Replica set name
	ReadPreference string `json:"readPreference,omitempty"` // Read preference
	ReadConcern    string `json:"readConcern,omitempty"`    // Read concern level
	WriteConcern   string `json:"writeConcern,omitempty"`   // Write concern level
	WTimeoutMS     int    `json:"wtimeoutMS,omitempty"`     // Write timeout in milliseconds
	AuthSource     string `json:"authSource,omitempty"`     // Authentication database (default: admin for root users)
	QueryLogging   bool   `json:"queryLogging,omitempty"`   // Enable MongoDB query logging
}

// Option is a functional option for configuring MongoDB connection.
// Use With* functions to create options (e.g., WithReplicaSet("rs0"), WithSRV(), etc.)
type Option func(*MongoOptions)
