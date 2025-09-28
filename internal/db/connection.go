package db

import (
	"context"
	"errors"
	"time"

	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrInvalidConnURL      = errors.New("failed to connect to DB, as the connection string is invalid")
	ErrConnectionEstablish = errors.New("failed to establish connection to DB")
	ErrClientInit          = errors.New("failed to initialize DB client")
	ErrConnectionLeak      = errors.New("unable to disconnect from DB, potential connection leak")
	ErrPingDB              = errors.New("failed to ping DB")
)

const (
	DefaultClientConnectTimeout = 10 * time.Second
)

type MongoDatabase interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

type MongoManager interface {
	Database() MongoDatabase
	Ping() error
	Disconnect() error
}

// MongoCredentials represents MongoDB authentication credentials
type MongoCredentials struct {
	Username string `json:"username,omitempty" log:"-"`
	Password string `json:"password,omitempty" log:"-"`
}

// ConnectionManager - Manages the connection to the underlying database.
type ConnectionManager struct {
	connectionURL string
	client        *mongo.Client
	database      *mongo.Database
	logger        *logger.AppLogger
	options       *MongoOptions
}

// NewMongoManager - Initializes DB connection and returns a Manager object which can be used to perform DB operations.
// The database parameter is optional - if empty, you can select databases later using DatabaseByName().
func NewMongoManager(hosts []string, database string, creds *MongoCredentials, lgr *logger.AppLogger, connOptions ...Option) (*ConnectionManager, error) {
	connURL, opts, err := ConnectionURL(hosts, database, creds, connOptions...)
	if err != nil {
		lgr.Error().Err(err).Msg("failed to build connection URL")
		return nil, err
	}

	lgr.Info().Str("connURL", MaskConnectionURL(connURL)).Msg("connecting to DB")

	connMgr := &ConnectionManager{
		logger:        lgr,
		connectionURL: connURL,
		options:       opts,
	}

	var c *mongo.Client
	if c, err = connMgr.newClient(); err == nil {
		var db *mongo.Database
		if database != "" {
			db = c.Database(database)
		}
		connMgr.database = db
		connMgr.client = c
		// Verify connection
		if pErr := connMgr.Ping(); pErr != nil {
			return nil, ErrConnectionEstablish
		}
		return connMgr, nil
	}
	return nil, err
}

// newClient - creates a new Mongo Client to connect DB.
func (c *ConnectionManager) newClient() (*mongo.Client, error) {
	var cmdMonitor *event.CommandMonitor
	if c.options.QueryLogging {
		cmdMonitor = &event.CommandMonitor{
			Started: func(_ context.Context, evt *event.CommandStartedEvent) {
				c.logger.Info().Str("dbQuery", evt.Command.String()).Send()
			},
		}
	}
	clientOptions := options.Client().ApplyURI(c.connectionURL).SetMonitor(cmdMonitor)
	ctx, cancel := context.WithTimeout(context.Background(), DefaultClientConnectTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to create new client")
		return nil, ErrClientInit
	}

	return client, nil
}

// Database - Returns configured database instance.
// If no database was specified during connection, this will return nil.
// Use DatabaseByName() to get a specific database at runtime.
func (c *ConnectionManager) Database() MongoDatabase {
	return c.database
}

// DatabaseByName - Returns a database instance for the specified name.
func (c *ConnectionManager) DatabaseByName(name string) MongoDatabase {
	return c.client.Database(name)
}

// Ping - Validates application's connectivity to the underlying database by pinging.
func (c *ConnectionManager) Ping() error {
	if err := c.client.Ping(context.TODO(), readpref.Primary()); err != nil {
		c.logger.Error().Err(err).Msg("failed to ping DB")
		return ErrPingDB
	}
	return nil
}

// Disconnect - Close connection to Database.
func (c *ConnectionManager) Disconnect() error {
	if err := c.client.Disconnect(context.Background()); err != nil {
		c.logger.Error().Err(err).Msg("failed to disconnect from DB")
		return ErrConnectionLeak
	}
	c.logger.Info().Msg("successfully disconnected from DB")
	return nil
}
