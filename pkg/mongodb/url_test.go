package mongodb_test

import (
	"testing"

	"github.com/rameshsunkara/go-rest-api-example/pkg/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectionURL(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	result := mongodb.MaskConnectionURL("mongodb://user:pass@localhost:27017/db")
	assert.Equal(t, "mongodb://%2A%2A%2A:%2A%2A%2A@localhost:27017/db", result)
}

func TestMaskConnectionURLNoCredentials(t *testing.T) {
	t.Parallel()
	result := mongodb.MaskConnectionURL("mongodb://localhost:27017/db")
	assert.Equal(t, "mongodb://localhost:27017/db", result)
}

func TestCredentialFromSideCarDefault(t *testing.T) {
	t.Parallel()
	// Test that function handles missing file gracefully
	_, err := mongodb.CredentialFromSideCar("")
	assert.Error(t, err)
	assert.Equal(t, mongodb.ErrSideCarFileRead, err)
}

func TestCredentialFromSideCarCustomFile(t *testing.T) {
	t.Parallel()
	// Test that function handles missing custom file gracefully
	_, err := mongodb.CredentialFromSideCar("/nonexistent/file.json")
	assert.Error(t, err)
	assert.Equal(t, mongodb.ErrSideCarFileRead, err)
}

func TestConstants(t *testing.T) {
	t.Parallel()
	// Test that MongoDB constants are properly defined
	assert.Equal(t, "mongodb", mongodb.MongoScheme)
	assert.Equal(t, "mongodb+srv", mongodb.MongoSRVScheme)
	assert.Equal(t, "/secrets/db.json", mongodb.DefaultMongoDBSidecar)
	assert.Equal(t, 5000, mongodb.DefaultWriteTimeoutMS)
}

func TestErrors(t *testing.T) {
	t.Parallel()
	// Test that error variables are properly defined
	assert.NotNil(t, mongodb.ErrNoHosts)
	assert.NotNil(t, mongodb.ErrSRVRequiresOneHost)
	assert.NotNil(t, mongodb.ErrInvalidReadPref)
	assert.NotNil(t, mongodb.ErrInvalidReadConcern)
	assert.NotNil(t, mongodb.ErrInvalidWriteConcern)
	assert.NotNil(t, mongodb.ErrSideCarFileRead)
	assert.NotNil(t, mongodb.ErrSideCarFileFormat)
}
