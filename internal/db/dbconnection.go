package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
	"time"

	"github.com/rameshsunkara/go-rest-api-example/internal/config"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var connectOnce sync.Once

// ConnectionTimeOut - Max time to establish DB connection // TODO: Move to config
const ConnectionTimeOut = 10 * time.Second

// Manager - Database manager
type Manager struct {
	client   *mongo.Client
	database *mongo.Database
}

// Init - Initializes DB connection and returns a Manager object which can be used to perform DB operations
func Init(dbName string) (*Manager, error) {
	c := config.GetConfig()
	connUrl := c.GetString("db.dsn")
	log.Debug().Str("DB Connection Url", connUrl)

	dbMgr := &Manager{}
	var connErr error
	connectOnce.Do(func() {
		c, err := newConnection(connUrl)
		if err != nil {
			connErr = err
		} else {
			db := c.Database(dbName)
			dbMgr.database = db
			dbMgr.client = c
		}
	})

	return dbMgr, connErr
}

func newConnection(connectionUrl string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(connectionUrl)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Error().Err(err).Msg("Connection Failed to Database")
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeOut)
	defer cancel()
	connErr := client.Connect(ctx)
	if connErr != nil {
		log.Error().Err(connErr).Msg("Connection Failed to Database")
		return nil, err
	}

	return client, nil
}

func (d *Manager) Client() (*mongo.Client, error) {
	if d.client == nil {
		return nil, errors.New("invalid state, database.Init is not called")
	}
	return d.client, nil
}

func (d *Manager) Database() (*mongo.Database, error) {
	if d.database == nil {
		return nil, errors.New("invalid state, database.Init is not called")
	}
	return d.database, nil
}

func (d *Manager) Ping() error {
	err := d.client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to DB")
	}

	return err
}
