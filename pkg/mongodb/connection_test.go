package mongodb_test

import (
	"testing"
	"time"

	"github.com/bogdanutanu/go-rest-api-example/pkg/mongodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testHost = "localhost:27017"

func TestConnectionURLValidation(t *testing.T) {
	creds := &mongodb.MongoCredentials{
		Username: "testuser",
		Password: "testpass",
	}

	t.Run("successful URL generation", func(t *testing.T) {
		url, opts, err := mongodb.ConnectionURL(testHost, "testdb", creds,
			mongodb.WithAuthSource("admin"),
			mongodb.WithQueryLogging(true))

		require.NoError(t, err)
		assert.NotEmpty(t, url)
		assert.NotNil(t, opts)
		assert.Contains(t, url, "testuser")
		assert.Contains(t, url, "testdb")
		assert.Contains(t, url, "authSource=admin")
	})
}

func TestConnectionURLInvalidHosts(t *testing.T) {
	creds := &mongodb.MongoCredentials{
		Username: "testuser",
		Password: "testpass",
	}

	tests := []struct {
		name     string
		hosts    string
		database string
		wantErr  error
	}{
		{
			name:     "empty hosts",
			hosts:    "",
			database: "testdb",
			wantErr:  mongodb.ErrNoHosts,
		},
		{
			name:     "whitespace only hosts",
			hosts:    "   ",
			database: "testdb",
			wantErr:  mongodb.ErrNoHosts,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := mongodb.ConnectionURL(tt.hosts, tt.database, creds)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestConnectionURLInvalidOptions(t *testing.T) {
	creds := &mongodb.MongoCredentials{
		Username: "testuser",
		Password: "testpass",
	}

	tests := []struct {
		name     string
		hosts    string
		database string
		options  []mongodb.Option
		wantErr  error
	}{
		{
			name:     "SRV with multiple hosts",
			hosts:    "host1:27017,host2:27017",
			database: "testdb",
			options:  []mongodb.Option{mongodb.WithSRV()},
			wantErr:  mongodb.ErrSRVRequiresOneHost,
		},
		{
			name:     "invalid read preference",
			hosts:    testHost,
			database: "testdb",
			options:  []mongodb.Option{mongodb.WithReadPreference("invalid")},
			wantErr:  mongodb.ErrInvalidReadPref,
		},
		{
			name:     "invalid read concern",
			hosts:    testHost,
			database: "testdb",
			options:  []mongodb.Option{mongodb.WithReadConcern("invalid")},
			wantErr:  mongodb.ErrInvalidReadConcern,
		},
		{
			name:     "invalid write concern",
			hosts:    testHost,
			database: "testdb",
			options:  []mongodb.Option{mongodb.WithWriteConcern("invalid")},
			wantErr:  mongodb.ErrInvalidWriteConcern,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := mongodb.ConnectionURL(tt.hosts, tt.database, creds, tt.options...)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestConnectionManagerStruct(t *testing.T) {
	t.Run("DatabaseByName with nil client panics as expected", func(t *testing.T) {
		connMgr := &mongodb.ConnectionManager{}

		// This will panic in real usage with nil client, which is expected behavior
		assert.Panics(t, func() {
			connMgr.DatabaseByName("test")
		})
	})
}

func TestConnectionManagerOptions(t *testing.T) {
	t.Run("basic functionality test", func(t *testing.T) {
		connMgr := &mongodb.ConnectionManager{}
		assert.NotNil(t, connMgr)
	})
}

func TestDefaultClientConnectTimeout(t *testing.T) {
	expectedTimeout := 10 * time.Second
	assert.Equal(t, expectedTimeout, mongodb.DefaultClientConnectTimeout)
}

func TestConnectionManagerErrors(t *testing.T) {
	tests := []struct {
		name         string
		expectedErr  error
		expectedText string
	}{
		{
			name:         "ErrInvalidConnURL",
			expectedErr:  mongodb.ErrInvalidConnURL,
			expectedText: "failed to connect to DB, as the connection string is invalid",
		},
		{
			name:         "mongodb.ErrConnectionEstablish",
			expectedErr:  mongodb.ErrConnectionEstablish,
			expectedText: "failed to establish connection to DB",
		},
		{
			name:         "ErrClientInit",
			expectedErr:  mongodb.ErrClientInit,
			expectedText: "failed to initialize DB client",
		},
		{
			name:         "ErrConnectionLeak",
			expectedErr:  mongodb.ErrConnectionLeak,
			expectedText: "unable to disconnect from DB, potential connection leak",
		},
		{
			name:         "ErrPingDB",
			expectedErr:  mongodb.ErrPingDB,
			expectedText: "failed to ping DB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedText, tt.expectedErr.Error())
			// Error is valid as expected
			assert.Error(t, tt.expectedErr)
		})
	}
}
