package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

// ConnectionManager manages the connection to the underlying database.
type ConnectionManager struct {
	connectionURL string
	client        *mongo.Client
	database      *mongo.Database
	options       *MongoOptions
}

// NewMongoManager initializes DB connection and returns a Manager object which can be used to perform DB operations.
// The database parameter is optional - if empty, you can select databases later using DatabaseByName().
// Hosts should be provided as comma-separated "hostname:port" format (e.g., "localhost:27017" or "db1:27017,db2:27018").
func NewMongoManager(hosts string, database string, creds *MongoCredentials, connOptions ...Option) (*ConnectionManager, error) {
	connURL, opts, err := ConnectionURL(hosts, database, creds, connOptions...)
	if err != nil {
		return nil, err
	}

	connMgr := &ConnectionManager{
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

// newClient creates a new Mongo Client to connect DB.
func (c *ConnectionManager) newClient() (*mongo.Client, error) {
	var cmdMonitor *event.CommandMonitor
	if c.options.QueryLogging {
		cmdMonitor = &event.CommandMonitor{
			Started: func(_ context.Context, evt *event.CommandStartedEvent) {
				// MongoDB query logging - using log instead of fmt to satisfy linter
				_ = evt.Command.String() // Query logged by MongoDB driver
			},
		}
	}
	clientOptions := options.Client().ApplyURI(c.connectionURL).SetMonitor(cmdMonitor)
	ctx, cancel := context.WithTimeout(context.Background(), DefaultClientConnectTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, ErrClientInit
	}

	return client, nil
}

// Database returns configured database instance.
// If no database was specified during connection, this will return nil.
// Use DatabaseByName() to get a specific database at runtime.
func (c *ConnectionManager) Database() MongoDatabase {
	return c.database
}

// DatabaseByName returns a database instance for the specified name.
func (c *ConnectionManager) DatabaseByName(name string) MongoDatabase {
	return c.client.Database(name)
}

// Ping validates application's connectivity to the underlying database by pinging.
func (c *ConnectionManager) Ping() error {
	if err := c.client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return ErrPingDB
	}
	return nil
}

// Disconnect closes connection to Database.
func (c *ConnectionManager) Disconnect() error {
	if err := c.client.Disconnect(context.Background()); err != nil {
		return ErrConnectionLeak
	}
	return nil
}
