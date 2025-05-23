package db_test

import (
	"testing"

	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/rameshsunkara/go-rest-api-example/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testLgr = logger.Setup("info", "test")
)

func TestNewMongoManager_InvalidConnURL(t *testing.T) {
	t.Parallel()
	creds := &db.MongoDBCredentials{}

	d, dErr := db.NewMongoManager(creds, nil, testLgr)
	assert.Nil(t, d)
	require.Error(t, dErr)
	assert.Equal(t, db.ErrInvalidConnURL, dErr)
}

func TestNewMongoManager_InvalidClient(t *testing.T) {
	t.Parallel()
	creds := &db.MongoDBCredentials{
		Hostname: "non-existent-hostname",
	}

	d, dErr := db.NewMongoManager(creds, nil, testLgr)
	assert.Nil(t, d)
	require.Error(t, dErr)
	assert.Equal(t, db.ErrConnectionEstablish, dErr)
}

func TestFillConnectionOpts(t *testing.T) {
	t.Parallel()
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
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			actual := db.FillConnectionOpts(tc.input)
			assert.Equal(t, tc.output, *actual, "test case %d:%s failed", i, tc.description)
		})
	}
}
