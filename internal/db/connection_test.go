package db_test

import (
	"context"
	"os"
	"runtime"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/log"
	"github.com/rameshsunkara/go-rest-api-example/internal/types"
	"github.com/rameshsunkara/go-rest-api-example/internal/util"
	"github.com/rameshsunkara/strikememongo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	testDBMgr db.MongoManager
)

const AppleChip = "arm64"

func mongoOptions() *strikememongo.Options {
	mongoVersion := "6.0.5"

	downloadUrl := ""
	if runtime.GOARCH == AppleChip {
		downloadUrl = "https://fastdl.mongodb.org/osx/mongodb-macos-arm64-6.0.5.tgz"
	}

	opts := &strikememongo.Options{
		MongoVersion: mongoVersion,
		DownloadURL:  downloadUrl,
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
		Hostname: mongoServer.URI(),
	}
	logger := log.Setup("test")
	d, dErr := db.NewMongoManager(creds, nil, logger)
	if dErr != nil {
		logger.Fatal().Err(dErr)
	}
	defer func(d *db.ConnectionManager) {
		err := d.Disconnect()
		if err != nil {
			logger.Error().Err(err).Msg("unable to disconnect from db")
		}
	}(d)
	testDBMgr = d
	insertTestData(logger)

	os.Exit(m.Run())
}

func insertTestData(logger *log.AppLogger) {
	database := testDBMgr.Database()
	dSvc := db.NewOrdersRepo(database)
	for i := 0; i < 500; i++ {
		product := []types.Product{
			{
				Name:        faker.Name(),
				Price:       util.RandomPrice(),
				Description: faker.Sentence(),
				UpdatedAt:   faker.TimeString(),
			},
			{
				Name:        faker.Name(),
				Price:       util.RandomPrice(),
				Description: faker.Sentence(),
				UpdatedAt:   faker.TimeString(),
			},
		}

		po := &types.Order{
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
	assert.Nil(t, err)
}
