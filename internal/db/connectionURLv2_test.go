package db_test

import (
	"errors"
	"testing"

	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectionURL(t *testing.T) {
	type connectionURLTestCase struct {
		Description string
		Hosts       []string
		Options     []db.Option
		ExpectedOut string
		ExpectedErr error
	}

	testCases := []connectionURLTestCase{
		{
			Description: "ensure ConnectionURL returns error when no hosts provided",
			Hosts:       []string{},
			Options:     []db.Option{},
			ExpectedOut: "",
			ExpectedErr: db.ErrNoHosts,
		},
		{
			Description: "basic connection URL with single host",
			Hosts:       []string{"localhost"},
			Options:     []db.Option{},
			ExpectedOut: "mongodb://localhost:27017",
		},
		{
			Description: "connection URL with custom port",
			Hosts:       []string{"localhost"},
			Options:     []db.Option{db.WithPort(27018)},
			ExpectedOut: "mongodb://localhost:27018",
		},
		{
			Description: "connection URL with credentials",
			Hosts:       []string{"localhost"},
			Options:     []db.Option{db.WithCredentials("user", "pass")},
			ExpectedOut: "mongodb://user:pass@localhost:27017",
		},
		{
			Description: "connection URL with database",
			Hosts:       []string{"localhost"},
			Options:     []db.Option{db.WithDatabase("mydb")},
			ExpectedOut: "mongodb://localhost:27017/mydb",
		},
		{
			Description: "connection URL with replica set",
			Hosts:       []string{"mongo1", "mongo2", "mongo3"},
			Options:     []db.Option{db.WithReplicaSet("rs0")},
			ExpectedOut: "mongodb://mongo1:27017,mongo2:27017,mongo3:27017?replicaSet=rs0",
		},
		{
			Description: "SRV connection URL",
			Hosts:       []string{"cluster.mongodb.net"},
			Options:     []db.Option{db.WithSRV()},
			ExpectedOut: "mongodb+srv://cluster.mongodb.net",
		},
		{
			Description: "connection URL with read preference",
			Hosts:       []string{"localhost"},
			Options:     []db.Option{db.WithReadPreference("secondary")},
			ExpectedOut: "mongodb://localhost:27017?readPreference=secondary",
		},
		{
			Description: "connection URL with write concern",
			Hosts:       []string{"localhost"},
			Options:     []db.Option{db.WithWriteConcern("majority")},
			ExpectedOut: "mongodb://localhost:27017?w=majority",
		},
		{
			Description: "connection URL with read concern",
			Hosts:       []string{"localhost"},
			Options:     []db.Option{db.WithReadConcern("majority")},
			ExpectedOut: "mongodb://localhost:27017?readConcernLevel=majority",
		},
		{
			Description: "connection URL with write timeout",
			Hosts:       []string{"localhost"},
			Options:     []db.Option{db.WithWriteTimeout(5000)},
			ExpectedOut: "mongodb://localhost:27017?wtimeoutMS=5000",
		},
		{
			Description: "comprehensive connection URL with all options",
			Hosts:       []string{"mongo1", "mongo2"},
			Options: []db.Option{
				db.WithCredentials("admin", "secret"),
				db.WithDatabase("production"),
				db.WithReplicaSet("rs0"),
				db.WithReadPreference("secondaryPreferred"),
				db.WithWriteConcern("majority"),
				db.WithReadConcern("majority"),
				db.WithWriteTimeout(10000),
			},
			ExpectedOut: "mongodb://admin:secret@mongo1:27017,mongo2:27017/production?readConcernLevel=majority&readPreference=secondaryPreferred&replicaSet=rs0&w=majority&wtimeoutMS=10000",
		},
		{
			Description: "host with port already specified",
			Hosts:       []string{"localhost:27018"},
			Options:     []db.Option{},
			ExpectedOut: "mongodb://localhost:27018",
		},
		{
			Description: "mixed hosts with and without ports",
			Hosts:       []string{"mongo1:27018", "mongo2"},
			Options:     []db.Option{db.WithPort(27019)},
			ExpectedOut: "mongodb://mongo1:27018,mongo2:27019",
		},
		{
			Description: "URL encoding in credentials",
			Hosts:       []string{"localhost"},
			Options:     []db.Option{db.WithCredentials("user@domain", "p@ss:w0rd")},
			ExpectedOut: "mongodb://user%2540domain:p%2540ss%253Aw0rd@localhost:27017",
		},
	}

	for i, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			result, err := db.ConnectionURL(tc.Hosts, tc.Options...)

			if tc.ExpectedErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tc.ExpectedErr), "Expected error %v, got %v", tc.ExpectedErr, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.ExpectedOut, result, "Test Case %d: %s", i+1, tc.Description)
			}
		})
	}
}

func TestConnectionURLValidationErrors(t *testing.T) {
	type validationTestCase struct {
		Description string
		Hosts       []string
		Options     []db.Option
		ExpectedErr error
	}

	testCases := []validationTestCase{
		{
			Description: "SRV with multiple hosts should return error",
			Hosts:       []string{"host1", "host2"},
			Options:     []db.Option{db.WithSRV()},
			ExpectedErr: db.ErrSRVRequiresOneHost,
		},
		{
			Description: "invalid port should return error",
			Hosts:       []string{"localhost"},
			Options:     []db.Option{db.WithPort(70000)},
			ExpectedErr: db.ErrInvalidPort,
		},
		{
			Description: "negative port should return error",
			Hosts:       []string{"localhost"},
			Options:     []db.Option{db.WithPort(-1)},
			ExpectedErr: db.ErrInvalidPort,
		},
		{
			Description: "invalid read preference should return error",
			Hosts:       []string{"localhost"},
			Options:     []db.Option{db.WithReadPreference("invalid")},
			ExpectedErr: db.ErrInvalidReadPref,
		},
		{
			Description: "invalid write concern should return error",
			Hosts:       []string{"localhost"},
			Options:     []db.Option{db.WithWriteConcern("invalid")},
			ExpectedErr: db.ErrInvalidWriteConcern,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			_, err := db.ConnectionURL(tc.Hosts, tc.Options...)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.ExpectedErr.Error())
		})
	}
}

func TestConnectionURLEdgeCases(t *testing.T) {
	t.Run("empty username with password", func(t *testing.T) {
		result, err := db.ConnectionURL([]string{"localhost"}, db.WithPassword("password"))
		require.NoError(t, err)
		assert.Equal(t, "mongodb://localhost:27017", result)
	})

	t.Run("username without password", func(t *testing.T) {
		result, err := db.ConnectionURL([]string{"localhost"}, db.WithUsername("user"))
		require.NoError(t, err)
		assert.Equal(t, "mongodb://user@localhost:27017", result)
	})

	t.Run("empty database name", func(t *testing.T) {
		result, err := db.ConnectionURL([]string{"localhost"}, db.WithDatabase(""))
		require.NoError(t, err)
		assert.Equal(t, "mongodb://localhost:27017", result)
	})

	t.Run("zero port uses default", func(t *testing.T) {
		result, err := db.ConnectionURL([]string{"localhost"}, db.WithPort(0))
		require.NoError(t, err)
		assert.Equal(t, "mongodb://localhost:27017", result)
	})
}

func TestMaskConnectionURL(t *testing.T) {
	testCases := []struct {
		Description string
		Input       string
		ExpectedOut string
	}{
		{
			Description: "mask credentials in connection URL",
			Input:       "mongodb://user:password@localhost:27017/mydb",
			ExpectedOut: "mongodb://%2A%2A%2A:%2A%2A%2A@localhost:27017/mydb",
		},
		{
			Description: "handle URL without credentials",
			Input:       "mongodb://localhost:27017/mydb",
			ExpectedOut: "mongodb://localhost:27017/mydb",
		},
		{
			Description: "handle empty URL",
			Input:       "",
			ExpectedOut: "",
		},
		{
			Description: "handle invalid URL (returns original)",
			Input:       "invalid-url",
			ExpectedOut: "invalid-url",
		},
		{
			Description: "mask SRV connection URL",
			Input:       "mongodb+srv://user:password@cluster.mongodb.net/mydb?retryWrites=true",
			ExpectedOut: "mongodb+srv://%2A%2A%2A:%2A%2A%2A@cluster.mongodb.net/mydb?retryWrites=true",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			result := db.MaskConnectionURL(tc.Input)
			assert.Equal(t, tc.ExpectedOut, result)
		})
	}
}
