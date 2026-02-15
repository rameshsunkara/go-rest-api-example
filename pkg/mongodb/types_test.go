package mongodb_test

import (
	"testing"

	"github.com/bogdanutanu/go-rest-api-example/pkg/mongodb"
	"github.com/stretchr/testify/assert"
)

func TestMongoCredentialsEmpty(t *testing.T) {
	t.Parallel()
	creds := mongodb.MongoCredentials{}

	assert.Empty(t, creds.Username)
	assert.Empty(t, creds.Password)
}

func TestMongoCredentialsFilled(t *testing.T) {
	t.Parallel()
	creds := mongodb.MongoCredentials{
		Username: "testuser",
		Password: "testpass",
	}

	assert.Equal(t, "testuser", creds.Username)
	assert.Equal(t, "testpass", creds.Password)
}

func TestMongoOptions(t *testing.T) {
	t.Parallel()
	opts := mongodb.MongoOptions{
		UseSRV:         true,
		ReplicaSet:     "rs0",
		ReadPreference: "secondary",
		ReadConcern:    "majority",
		WriteConcern:   "majority",
		WTimeoutMS:     5000,
		AuthSource:     "admin",
		QueryLogging:   true,
	}

	assert.True(t, opts.UseSRV)
	assert.Equal(t, "rs0", opts.ReplicaSet)
	assert.Equal(t, "secondary", opts.ReadPreference)
	assert.Equal(t, "majority", opts.ReadConcern)
	assert.Equal(t, "majority", opts.WriteConcern)
	assert.Equal(t, 5000, opts.WTimeoutMS)
	assert.Equal(t, "admin", opts.AuthSource)
	assert.True(t, opts.QueryLogging)
}

func TestConnectionManager(t *testing.T) {
	t.Parallel()
	connMgr := &mongodb.ConnectionManager{}
	assert.NotNil(t, connMgr)
	assert.Nil(t, connMgr.Database())
}
