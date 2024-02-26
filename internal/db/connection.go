package db

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrInvalidConnURL      = errors.New("failed to connect to DB, as the connection string is invalid")
	ErrConnectionEstablish = errors.New("failed to establish connection to DB")
	ErrClientInit          = errors.New("failed to initialize db client")
	ErrConnectionLeak      = errors.New("unable to disconnect from db, potential connection leak")
	ErrPingDB              = errors.New("failed to ping DB")
)

const (
	DefaultConnTimeout = 10 * time.Second
	DefaultDatabase    = "ecommerce"
)

type MongoDatabase interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

type MongoManager interface {
	Database() MongoDatabase
	Ping() error
	Disconnect() error
}

type ConnectionOpts struct {
	ConnectionTimeout time.Duration
	PrintQueries      bool
	Database          string
}

// ConnectionManager - Manages the connection to the underlying database.
type ConnectionManager struct {
	connectionURL string
	client        *mongo.Client
	database      *mongo.Database
	credentials   MongoDBCredentials
	logger        zerolog.Logger
}

// NewMongoManager - Initializes DB connection and returns a Manager object which can be used to perform DB operations.
func NewMongoManager(mc MongoDBCredentials, opts *ConnectionOpts, log zerolog.Logger) (*ConnectionManager, error) {
	connURL := MongoConnectionURL(mc)
	if len(connURL) == 0 {
		return nil, ErrInvalidConnURL
	}
	connMgr := &ConnectionManager{
		credentials:   mc,
		logger:        log,
		connectionURL: connURL,
	}
	connOpts := fillConnectionOpts(opts)
	var err error
	var c *mongo.Client
	if c, err = connMgr.newClient(connOpts); err == nil {
		db := c.Database(opts.Database)
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

func fillConnectionOpts(opts *ConnectionOpts) *ConnectionOpts {
	if opts == nil {
		return &ConnectionOpts{
			PrintQueries:      false,
			ConnectionTimeout: DefaultConnTimeout,
			Database:          DefaultDatabase,
		}
	}
	if opts.ConnectionTimeout == 0 {
		opts.ConnectionTimeout = DefaultConnTimeout
	}
	if opts.Database == "" {
		opts.Database = DefaultDatabase
	}
	return opts
}

// newClient - creates a new Mongo Client to connect DB.
func (c *ConnectionManager) newClient(connOpts *ConnectionOpts) (*mongo.Client, error) {
	var cmdMonitor *event.CommandMonitor
	if connOpts.PrintQueries {
		cmdMonitor = &event.CommandMonitor{
			Started: func(_ context.Context, evt *event.CommandStartedEvent) {
				c.logger.Debug().Interface("dbQuery", evt.Command).Send()
			},
		}
	}
	clientOptions := options.Client().ApplyURI(c.connectionURL).SetMonitor(cmdMonitor)
	ctx, cancel := context.WithTimeout(context.Background(), connOpts.ConnectionTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		c.logger.Err(err).Msg("failed to create new client")
		return nil, ErrClientInit
	}

	return client, nil
}

// Database - Returns configured database instance.
func (c *ConnectionManager) Database() MongoDatabase {
	return c.database
}

// Ping - Validates application's connectivity to the underlying database by pinging.
func (c *ConnectionManager) Ping() error {
	if err := c.client.Ping(context.TODO(), readpref.Primary()); err != nil {
		c.logger.Err(err).Msg("failed to ping DB")
		return ErrPingDB
	}
	return nil
}

// Disconnect - Close connection to Database.
func (c *ConnectionManager) Disconnect() error {
	if err := c.client.Disconnect(context.Background()); err != nil {
		c.logger.Err(err).Msg("failed to disconnect from DB")
		return ErrConnectionLeak
	}
	c.logger.Info().Msg("successfully disconnected from DB")
	return nil
}
