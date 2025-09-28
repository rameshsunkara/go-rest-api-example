package mongodb_test

import (
	"testing"

	"github.com/rameshsunkara/go-rest-api-example/pkg/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectionURL(t *testing.T) {
	creds := &mongodb.MongoCredentials{
		Username: "testuser",
		Password: "testpass",
	}

	url, opts, err := mongodb.ConnectionURL("localhost:27017", "testdb", creds)

	require.NoError(t, err)
	require.NotNil(t, opts)
	assert.Contains(t, url, "mongodb://")
}

func TestMaskConnectionURL(t *testing.T) {
	result := mongodb.MaskConnectionURL("mongodb://user:pass@localhost:27017/db")
	assert.Equal(t, "mongodb://%2A%2A%2A:%2A%2A%2A@localhost:27017/db", result)
}
