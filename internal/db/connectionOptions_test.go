package db_test

import (
	"testing"

	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestWithPort(t *testing.T) {
	opts := &db.MongoOptions{}
	option := db.WithPort(27018)
	option(opts)
	
	assert.Equal(t, 27018, opts.Port)
}

func TestWithSRV(t *testing.T) {
	opts := &db.MongoOptions{}
	option := db.WithSRV()
	option(opts)
	
	assert.True(t, opts.UseSRV)
}

func TestWithDatabase(t *testing.T) {
	opts := &db.MongoOptions{}
	option := db.WithDatabase("testdb")
	option(opts)
	
	assert.Equal(t, "testdb", opts.Database)
}

func TestWithUsername(t *testing.T) {
	opts := &db.MongoOptions{}
	option := db.WithUsername("testuser")
	option(opts)
	
	assert.Equal(t, "testuser", opts.Username)
}

func TestWithPassword(t *testing.T) {
	opts := &db.MongoOptions{}
	option := db.WithPassword("testpass")
	option(opts)
	
	assert.Equal(t, "testpass", opts.Password)
}

func TestWithCredentials(t *testing.T) {
	opts := &db.MongoOptions{}
	option := db.WithCredentials("user", "pass")
	option(opts)
	
	assert.Equal(t, "user", opts.Username)
	assert.Equal(t, "pass", opts.Password)
}

func TestWithReplicaSet(t *testing.T) {
	opts := &db.MongoOptions{}
	option := db.WithReplicaSet("rs0")
	option(opts)
	
	assert.Equal(t, "rs0", opts.ReplicaSet)
}

func TestWithReadPreference(t *testing.T) {
	opts := &db.MongoOptions{}
	option := db.WithReadPreference("secondary")
	option(opts)
	
	assert.Equal(t, "secondary", opts.ReadPreference)
}

func TestWithReadConcern(t *testing.T) {
	opts := &db.MongoOptions{}
	option := db.WithReadConcern("majority")
	option(opts)
	
	assert.Equal(t, "majority", opts.ReadConcern)
}

func TestWithWriteConcern(t *testing.T) {
	opts := &db.MongoOptions{}
	option := db.WithWriteConcern("majority")
	option(opts)
	
	assert.Equal(t, "majority", opts.WriteConcern)
}

func TestWithWriteTimeout(t *testing.T) {
	opts := &db.MongoOptions{}
	option := db.WithWriteTimeout(5000)
	option(opts)
	
	assert.Equal(t, 5000, opts.WTimeoutMS)
}

func TestIsValidReadPreference(t *testing.T) {
	testCases := []struct {
		Description string
		Input       string
		Expected    bool
	}{
		{
			Description: "primary is valid",
			Input:       "primary",
			Expected:    true,
		},
		{
			Description: "primaryPreferred is valid",
			Input:       "primaryPreferred",
			Expected:    true,
		},
		{
			Description: "secondary is valid",
			Input:       "secondary",
			Expected:    true,
		},
		{
			Description: "secondaryPreferred is valid",
			Input:       "secondaryPreferred",
			Expected:    true,
		},
		{
			Description: "nearest is valid",
			Input:       "nearest",
			Expected:    true,
		},
		{
			Description: "invalid preference returns false",
			Input:       "invalid",
			Expected:    false,
		},
		{
			Description: "empty string returns false",
			Input:       "",
			Expected:    false,
		},
		{
			Description: "case sensitive - Primary returns false",
			Input:       "Primary",
			Expected:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			result := db.IsValidReadPreference(tc.Input)
			assert.Equal(t, tc.Expected, result)
		})
	}
}

func TestIsValidReadConcern(t *testing.T) {
	testCases := []struct {
		Description string
		Input       string
		Expected    bool
	}{
		{
			Description: "local is valid",
			Input:       "local",
			Expected:    true,
		},
		{
			Description: "available is valid",
			Input:       "available",
			Expected:    true,
		},
		{
			Description: "majority is valid",
			Input:       "majority",
			Expected:    true,
		},
		{
			Description: "linearizable is valid",
			Input:       "linearizable",
			Expected:    true,
		},
		{
			Description: "snapshot is valid",
			Input:       "snapshot",
			Expected:    true,
		},
		{
			Description: "invalid concern returns false",
			Input:       "invalid",
			Expected:    false,
		},
		{
			Description: "empty string returns false",
			Input:       "",
			Expected:    false,
		},
		{
			Description: "case sensitive - Local returns false",
			Input:       "Local",
			Expected:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			result := db.IsValidReadConcern(tc.Input)
			assert.Equal(t, tc.Expected, result)
		})
	}
}

func TestIsValidWriteConcern(t *testing.T) {
	testCases := []struct {
		Description string
		Input       string
		Expected    bool
	}{
		{
			Description: "majority is valid",
			Input:       "majority",
			Expected:    true,
		},
		{
			Description: "0 is valid",
			Input:       "0",
			Expected:    true,
		},
		{
			Description: "1 is valid",
			Input:       "1",
			Expected:    true,
		},
		{
			Description: "2 is valid",
			Input:       "2",
			Expected:    true,
		},
		{
			Description: "3 is valid",
			Input:       "3",
			Expected:    true,
		},
		{
			Description: "invalid concern returns false",
			Input:       "invalid",
			Expected:    false,
		},
		{
			Description: "empty string returns false",
			Input:       "",
			Expected:    false,
		},
		{
			Description: "4 is invalid (not in our list)",
			Input:       "4",
			Expected:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			result := db.IsValidWriteConcern(tc.Input)
			assert.Equal(t, tc.Expected, result)
		})
	}
}

func TestMultipleOptions(t *testing.T) {
	t.Run("multiple options can be applied together", func(t *testing.T) {
		opts := &db.MongoOptions{}
		
		// Apply multiple options
		options := []db.Option{
			db.WithPort(27018),
			db.WithDatabase("testdb"),
			db.WithCredentials("user", "pass"),
			db.WithReplicaSet("rs0"),
			db.WithReadPreference("secondary"),
			db.WithWriteConcern("majority"),
			db.WithWriteTimeout(5000),
		}
		
		for _, opt := range options {
			opt(opts)
		}
		
		assert.Equal(t, 27018, opts.Port)
		assert.Equal(t, "testdb", opts.Database)
		assert.Equal(t, "user", opts.Username)
		assert.Equal(t, "pass", opts.Password)
		assert.Equal(t, "rs0", opts.ReplicaSet)
		assert.Equal(t, "secondary", opts.ReadPreference)
		assert.Equal(t, "majority", opts.WriteConcern)
		assert.Equal(t, 5000, opts.WTimeoutMS)
	})
	
	t.Run("options override previous values", func(t *testing.T) {
		opts := &db.MongoOptions{}
		
		// Set initial value
		db.WithPort(27017)(opts)
		assert.Equal(t, 27017, opts.Port)
		
		// Override with new value
		db.WithPort(27018)(opts)
		assert.Equal(t, 27018, opts.Port)
	})
}

func TestMongoOptionsZeroValues(t *testing.T) {
	opts := &db.MongoOptions{}
	
	// All fields should have zero values initially
	assert.Equal(t, 0, opts.Port)
	assert.False(t, opts.UseSRV)
	assert.Equal(t, "", opts.Database)
	assert.Equal(t, "", opts.Username)
	assert.Equal(t, "", opts.Password)
	assert.Equal(t, "", opts.ReplicaSet)
	assert.Equal(t, "", opts.ReadPreference)
	assert.Equal(t, "", opts.ReadConcern)
	assert.Equal(t, "", opts.WriteConcern)
	assert.Equal(t, 0, opts.WTimeoutMS)
}