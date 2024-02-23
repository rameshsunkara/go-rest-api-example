package db

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrInvalidConnUrl = errors.New("failed to connect to DB, as the connection string is invalid")
	ErrClientCreation = errors.New("failed to create new client to connect with db")
	ErrClientInit     = errors.New("failed to initialize db client")
	ErrConnectionLeak = errors.New("unable to disconnect from db, potential connection leak")
)

const DefConnectionTimeOut = 0 * time.Second

type MongoDatabase interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

type MongoManager interface {
	Database() MongoDatabase
	Ping() error
	Disconnect() error
}

// ConnectionManager - Manages the connection to the underlying database
type ConnectionManager struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoManager - Initializes DB connection and returns a Manager object which can be used to perform DB operations
func NewMongoManager(dbName string, connUrl string) (*ConnectionManager, error) {
	log.Info().Str("DB Connection Url", connUrl)

	dbMgr := &ConnectionManager{}
	if c, err := newClient(connUrl); err != nil {
		return nil, err
	} else {
		db := c.Database(dbName)
		dbMgr.database = db
		dbMgr.client = c
	}

	return dbMgr, nil
}

// newClient - creates a new Mongo Client to connect to the specified url and initializes the Client
func newClient(connectionUrl string) (*mongo.Client, error) {
	if len(connectionUrl) == 0 {
		return nil, ErrInvalidConnUrl
	}
	opts := options.Client().ApplyURI(connectionUrl)
	opts.SetConnectTimeout(DefConnectionTimeOut)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Error().Err(err).Msg("failed to establish db connection")
		return nil, ErrClientCreation
	}
	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Error().Err(err).Msg("failed to initialize db client")
		return nil, ErrClientInit
	}
	return client, nil
}

// Database - Returns configured database instance
func (c *ConnectionManager) Database() MongoDatabase {
	return c.database
}

// Ping - Validates application's connectivity to the underlying database by pinging
func (c *ConnectionManager) Ping() error {
	if err := c.client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Error().Err(err).Msg("unable to connect to DB")
		return err
	}
	return nil
}

// Disconnect - Close connection to Database
func (c *ConnectionManager) Disconnect() error {
	if err := c.client.Disconnect(context.Background()); err != nil {
		log.Error().Err(err).Msg("unable to disconnect from DB")
		return ErrConnectionLeak
	}
	log.Info().Msg("successfully disconnected from DB")
	return nil
}
