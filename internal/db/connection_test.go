package db_test

import (
	"context"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/rameshsunkara/go-rest-api-example/internal/models"
	"github.com/rameshsunkara/go-rest-api-example/internal/models/data"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
	"github.com/rameshsunkara/strikememongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	testDBMgr db.MongoManager
	lgr       = logger.Setup(models.ServiceEnv{Name: "test"})
)

const AppleChip = "arm64"

func mongoOptions() *strikememongo.Options {
	mongoVersion := "6.0.5"

	downloadURL := ""
	if runtime.GOARCH == AppleChip {
		downloadURL = "https://fastdl.mongodb.org/osx/mongodb-macos-arm64-6.0.5.tgz"
	}

	opts := &strikememongo.Options{
		MongoVersion: mongoVersion,
		DownloadURL:  downloadURL,
	}
	return opts
}

func TestMain(m *testing.M) {
	mongoServer, err := strikememongo.StartWithOptions(mongoOptions())
	if err != nil {
		panic(err)
	}
	defer mongoServer.Stop()
	creds := &db.MongoDBCredentials{
		Hostname: strings.TrimPrefix(mongoServer.URI(), "mongodb://"),
	}
	d, dErr := db.NewMongoManager(creds, nil, lgr)
	if dErr != nil {
		lgr.Fatal().Err(dErr)
	}
	defer func(d *db.ConnectionManager) {
		discErr := d.Disconnect()
		if discErr != nil {
			lgr.Error().Err(discErr).Msg("unable to disconnect from db")
		}
	}(d)
	testDBMgr = d
	insertTestData(lgr)
	m.Run()
}

func insertTestData(logger *logger.AppLogger) {
	database := testDBMgr.Database()
	dSvc := db.NewOrdersRepo(database, logger)
	for i := 0; i < 10; i++ {
		product := []data.Product{
			{
				Name:      faker.Name(),
				Price:     util.RandomPrice(),
				UpdatedAt: time.Now(),
			},
			{
				Name:      faker.Name(),
				Price:     util.RandomPrice(),
				UpdatedAt: time.Now(),
			},
		}

		po := &data.Order{
			Products: product,
		}
		_, err := dSvc.Create(context.TODO(), po)
		if err != nil {
			logger.Fatal().Err(err).Msg("unable to insert data")
		}
	}
}

func TestDatabase(t *testing.T) {
	d := testDBMgr.Database()
	assert.NotNil(t, d)
	assert.IsType(t, &mongo.Database{}, d)
}

func TestPing(t *testing.T) {
	err := testDBMgr.Ping()
	require.NoError(t, err)
}

func TestNewMongoManager_InvalidConnURL(t *testing.T) {
	creds := &db.MongoDBCredentials{}

	d, dErr := db.NewMongoManager(creds, nil, lgr)
	assert.Nil(t, d)
	require.Error(t, dErr)
	assert.EqualValues(t, db.ErrInvalidConnURL, dErr)
}

func TestNewMongoManager_InvalidClient(t *testing.T) {
	creds := &db.MongoDBCredentials{
		Hostname: "non-existent-hostname",
	}

	d, dErr := db.NewMongoManager(creds, nil, lgr)
	assert.Nil(t, d)
	require.Error(t, dErr)
	assert.EqualValues(t, db.ErrConnectionEstablish, dErr)
}

func TestFillConnectionOpts(t *testing.T) {
	testCases := []struct {
		description string
		input       *db.ConnectionOpts
		output      db.ConnectionOpts
	}{
		{
			description: "expect connect time out and database set to default",
			input: &db.ConnectionOpts{
				PrintQueries: true,
			},
			output: db.ConnectionOpts{
				Database:          db.DefDatabase,
				ConnectionTimeout: db.DefConnectionTimeOut,
				PrintQueries:      true,
			},
		},
		{
			description: "expect showQueries to be false",
			input: &db.ConnectionOpts{
				ConnectionTimeout: db.DefConnectionTimeOut,
			},
			output: db.ConnectionOpts{
				Database:          db.DefDatabase,
				ConnectionTimeout: db.DefConnectionTimeOut,
				PrintQueries:      false,
			},
		},
		{
			description: "expect connect time out set to default and showQueries to be false",
			input:       &db.ConnectionOpts{},
			output: db.ConnectionOpts{
				Database:          db.DefDatabase,
				ConnectionTimeout: db.DefConnectionTimeOut,
				PrintQueries:      false,
			},
		},
		{
			description: "expect connect time out set to default and showQueries to be false when input is nil",
			input:       nil,
			output: db.ConnectionOpts{
				Database:          db.DefDatabase,
				ConnectionTimeout: db.DefConnectionTimeOut,
				PrintQueries:      false,
			},
		},
	}

	for i, tc := range testCases {
		actual := db.FillConnectionOpts(tc.input)
		assert.Equal(t, tc.output, *actual, "test case %d:%s failed", i, tc.description)
	}
}
