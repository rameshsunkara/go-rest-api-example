package mongodb_test

import (
	"testing"

	"github.com/bogdanutanu/go-rest-api-example/pkg/mongodb"
	"github.com/stretchr/testify/assert"
)

func TestWithSRV(t *testing.T) {
	opts := &mongodb.MongoOptions{}
	option := mongodb.WithSRV()
	option(opts)

	assert.True(t, opts.UseSRV)
}

func TestWithReplicaSet(t *testing.T) {
	replicaSetName := "rs0"
	opts := &mongodb.MongoOptions{}
	option := mongodb.WithReplicaSet(replicaSetName)
	option(opts)

	assert.Equal(t, replicaSetName, opts.ReplicaSet)
}

func TestWithReadPreference(t *testing.T) {
	readPref := "secondary"
	opts := &mongodb.MongoOptions{}
	option := mongodb.WithReadPreference(readPref)
	option(opts)

	assert.Equal(t, readPref, opts.ReadPreference)
}

func TestWithReadConcern(t *testing.T) {
	readConcern := "majority"
	opts := &mongodb.MongoOptions{}
	option := mongodb.WithReadConcern(readConcern)
	option(opts)

	assert.Equal(t, readConcern, opts.ReadConcern)
}

func TestWithWriteConcern(t *testing.T) {
	writeConcern := "majority"
	opts := &mongodb.MongoOptions{}
	option := mongodb.WithWriteConcern(writeConcern)
	option(opts)

	assert.Equal(t, writeConcern, opts.WriteConcern)
}

func TestWithWriteTimeout(t *testing.T) {
	timeout := 5000
	opts := &mongodb.MongoOptions{}
	option := mongodb.WithWriteTimeout(timeout)
	option(opts)

	assert.Equal(t, timeout, opts.WTimeoutMS)
}

func TestWithAuthSource(t *testing.T) {
	authSource := "admin"
	opts := &mongodb.MongoOptions{}
	option := mongodb.WithAuthSource(authSource)
	option(opts)

	assert.Equal(t, authSource, opts.AuthSource)
}

func TestWithQueryLogging(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
	}{
		{"enable query logging", true},
		{"disable query logging", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &mongodb.MongoOptions{}
			option := mongodb.WithQueryLogging(tt.enabled)
			option(opts)

			assert.Equal(t, tt.enabled, opts.QueryLogging)
		})
	}
}

func TestIsValidReadPreference(t *testing.T) {
	tests := []struct {
		name     string
		pref     string
		expected bool
	}{
		{"primary", "primary", true},
		{"primaryPreferred", "primaryPreferred", true},
		{"secondary", "secondary", true},
		{"secondaryPreferred", "secondaryPreferred", true},
		{"nearest", "nearest", true},
		{"invalid preference", "invalid", false},
		{"empty preference", "", false},
		{"case sensitive", "PRIMARY", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mongodb.IsValidReadPreference(tt.pref)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidReadConcern(t *testing.T) {
	tests := []struct {
		name     string
		concern  string
		expected bool
	}{
		{"local", "local", true},
		{"available", "available", true},
		{"majority", "majority", true},
		{"linearizable", "linearizable", true},
		{"snapshot", "snapshot", true},
		{"invalid concern", "invalid", false},
		{"empty concern", "", false},
		{"case sensitive", "LOCAL", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mongodb.IsValidReadConcern(tt.concern)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidWriteConcern(t *testing.T) {
	tests := []struct {
		name     string
		concern  string
		expected bool
	}{
		{"majority", "majority", true},
		{"fire and forget", "0", true},
		{"primary only", "1", true},
		{"primary plus one", "2", true},
		{"primary plus two", "3", true},
		{"invalid concern", "invalid", false},
		{"empty concern", "", false},
		{"negative number", "-1", false},
		{"high number", "10", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mongodb.IsValidWriteConcern(tt.concern)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOptionsFunctionalPattern(t *testing.T) {
	// Test that multiple options can be applied together
	opts := &mongodb.MongoOptions{}

	options := []mongodb.Option{
		mongodb.WithSRV(),
		mongodb.WithReplicaSet("rs0"),
		mongodb.WithReadPreference("secondary"),
		mongodb.WithReadConcern("majority"),
		mongodb.WithWriteConcern("majority"),
		mongodb.WithWriteTimeout(10000),
		mongodb.WithAuthSource("admin"),
		mongodb.WithQueryLogging(true),
	}

	// Apply all options
	for _, option := range options {
		option(opts)
	}

	// Verify all options were applied
	assert.True(t, opts.UseSRV)
	assert.Equal(t, "rs0", opts.ReplicaSet)
	assert.Equal(t, "secondary", opts.ReadPreference)
	assert.Equal(t, "majority", opts.ReadConcern)
	assert.Equal(t, "majority", opts.WriteConcern)
	assert.Equal(t, 10000, opts.WTimeoutMS)
	assert.Equal(t, "admin", opts.AuthSource)
	assert.True(t, opts.QueryLogging)
}

func TestOptionsOverwrite(t *testing.T) {
	// Test that options can be overwritten
	opts := &mongodb.MongoOptions{}

	// Apply first set of options
	mongodb.WithReadPreference("primary")(opts)
	mongodb.WithQueryLogging(false)(opts)

	// Verify initial values
	assert.Equal(t, "primary", opts.ReadPreference)
	assert.False(t, opts.QueryLogging)

	// Overwrite with new values
	mongodb.WithReadPreference("secondary")(opts)
	mongodb.WithQueryLogging(true)(opts)

	// Verify new values
	assert.Equal(t, "secondary", opts.ReadPreference)
	assert.True(t, opts.QueryLogging)
}

func TestOptions(t *testing.T) {
	// Test behavior with zero values
	opts := &mongodb.MongoOptions{}

	mongodb.WithReplicaSet("")(opts)
	mongodb.WithReadPreference("")(opts)
	mongodb.WithWriteTimeout(0)(opts)

	assert.Empty(t, opts.ReplicaSet)
	assert.Empty(t, opts.ReadPreference)
	assert.Equal(t, 0, opts.WTimeoutMS)
}
