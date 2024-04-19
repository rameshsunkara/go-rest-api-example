package db_test

import (
	"context"
	"math/rand"
	"os"
	"runtime"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/types"
	"github.com/rameshsunkara/strikememongo"
	"github.com/rs/zerolog/log"
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

	d, dErr := db.NewMongoManager(strikememongo.RandomDatabase(), mongoServer.URI(), nil)
	if dErr != nil {
		log.Fatal().Err(dErr)
	}
	defer func(d *db.ConnectionManager) {
		err := d.Disconnect()
		if err != nil {
			log.Error().Err(err).Msg("unable to disconnect from db")
		}
	}(d)
	testDBMgr = d
	insertTestData()

	os.Exit(m.Run())
}

func insertTestData() {
	database := testDBMgr.Database()
	dSvc := db.NewOrdersRepo(database)

	for i := 0; i < 500; i++ {
		product := []types.Product{
			{
				Name:        faker.Name(),
				Price:       (uint)(rand.Intn(90) + 10),
				Description: faker.Sentence(),
				UpdatedAt:   faker.TimeString(),
			},
			{
				Name:        faker.Name(),
				Price:       (uint)(rand.Intn(1000) + 10),
				Description: faker.Sentence(),
				UpdatedAt:   faker.TimeString(),
			},
		}

		po := &types.Order{
			Products: product,
		}
		_, err := dSvc.Create(context.TODO(), po)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to insert data")
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
