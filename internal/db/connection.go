package db

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	InvalidConnUrlErr = errors.New("failed to connect to DB, as the connection string is invalid")
	ClientCreationErr = errors.New("failed to create new client to connect with db")
	ClientInitErr     = errors.New("failed to initialize db client")
	ConnectionLeak    = errors.New("unable to disconnect from db, potential connection leak")
	connectOnce       sync.Once
)

// ConnectionTimeOut - Max time to establish DB connection // TODO: Move to config
const ConnectionTimeOut = 10 * time.Second

type MongoDatabase interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

type MongoManager interface {
	Database() MongoDatabase
	Ping() error
	Disconnect() error
}

// connectionManager - Implements MongoManager
type connectionManager struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoManager - Initializes DB connection and returns a Manager object which can be used to perform DB operations
func NewMongoManager(dbName string, connUrl string) (MongoManager, error) {
	log.Debug().Str("DB Connection Url", connUrl)

	dbMgr := &connectionManager{}
	var connErr error
	connectOnce.Do(func() {
		if c, err := newClient(connUrl); err != nil {
			connErr = err
		} else {
			db := c.Database(dbName)
			dbMgr.database = db
			dbMgr.client = c

			// Verify connection
			if err := dbMgr.Ping(); err != nil {
				connErr = err
			}
		}
	})

	return dbMgr, connErr
}

// newClient - creates a new Mongo Client to connect to the specified url and initializes the Client
func newClient(connectionUrl string) (*mongo.Client, error) {
	if len(connectionUrl) == 0 {
		return nil, InvalidConnUrlErr
	}
	clientOptions := options.Client().ApplyURI(connectionUrl)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Error().Err(err).Msg("Connection Failed to Database")
		return nil, ClientCreationErr
	}

	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeOut)
	defer cancel()
	connErr := client.Connect(ctx)
	if connErr != nil {
		log.Error().Err(connErr).Msg("Connection Failed to Database")
		return nil, ClientInitErr
	}

	return client, nil
}

// Database - Returns configured database instance
func (c *connectionManager) Database() MongoDatabase {
	return c.database
}

// Ping - Validates application's connectivity to the underlying database by pinging
func (c *connectionManager) Ping() error {
	if err := c.client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Error().Err(err).Msg("unable to connect to DB")
		return err
	}
	return nil
}

// Disconnect - Close connection to Database
func (c *connectionManager) Disconnect() error {
	log.Info().Msg("Disconnecting from Database")
	if err := c.client.Disconnect(context.Background()); err != nil {
		log.Error().Err(err).Msg("unable to disconnect from DB")
		return ConnectionLeak
	}
	log.Info().Msg("Successfully disconnected from DB")
	return nil
}
