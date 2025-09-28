package db

// Valid MongoDB read preferences - controls WHERE to read from (server selection)
var validReadPreferences = []string{
	"primary",            // Read only from primary server
	"primaryPreferred",   // Read from primary, fall back to secondary
	"secondary",          // Read only from secondary servers
	"secondaryPreferred", // Read from secondary, fall back to primary
	"nearest",            // Read from server with lowest latency
}

// Valid MongoDB read concern levels - controls CONSISTENCY of data returned
var validReadConcerns = []string{
	"local",        // Return local data (may not be majority acknowledged)
	"available",    // Return immediately available data (fastest)
	"majority",     // Return data acknowledged by majority of replica set
	"linearizable", // Return data reflecting all successful majority writes
	"snapshot",     // Read from consistent snapshot (for transactions)
}

// Valid MongoDB write concern levels - controls ACKNOWLEDGMENT of write operations
var validWriteConcerns = []string{
	"majority", // Wait for acknowledgment from majority of replica set
	"0",        // No acknowledgment required (fire and forget)
	"1",        // Wait for acknowledgment from primary only
	"2",        // Wait for acknowledgment from primary + 1 secondary
	"3",        // Wait for acknowledgment from primary + 2 secondaries
}

// MongoOptions represents optional MongoDB connection settings
// Note: Hosts should be provided in "hostname:port" format directly to the ConnectionURL function
type MongoOptions struct {
	UseSRV         bool   `json:"useSRV,omitempty"`         // Use SRV connection
	ReplicaSet     string `json:"replicaSet,omitempty"`     // Replica set name
	ReadPreference string `json:"readPreference,omitempty"` // Read preference
	ReadConcern    string `json:"readConcern,omitempty"`    // Read concern level
	WriteConcern   string `json:"writeConcern,omitempty"`   // Write concern level
	WTimeoutMS     int    `json:"wtimeoutMS,omitempty"`     // Write timeout in milliseconds
	AuthSource     string `json:"authSource,omitempty"`     // Authentication database (default: admin for root users)
	QueryLogging   bool   `json:"queryLogging,omitempty"`   // Enable MongoDB query logging
}

// Option is a functional option for configuring MongoDB connection.
// Use With* functions to create options (e.g., WithPort(27017), WithSRV(), etc.)
type Option func(*MongoOptions)

// WithSRV enables SRV connection mode
func WithSRV() Option {
	return func(opts *MongoOptions) {
		opts.UseSRV = true
	}
}

// WithReplicaSet sets the replica set name
func WithReplicaSet(replicaSet string) Option {
	return func(opts *MongoOptions) {
		opts.ReplicaSet = replicaSet
	}
}

// WithReadPreference sets the read preference
func WithReadPreference(pref string) Option {
	return func(opts *MongoOptions) {
		opts.ReadPreference = pref
	}
}

// WithReadConcern sets the read concern
func WithReadConcern(concern string) Option {
	return func(opts *MongoOptions) {
		opts.ReadConcern = concern
	}
}

// WithWriteConcern sets the write concern
func WithWriteConcern(concern string) Option {
	return func(opts *MongoOptions) {
		opts.WriteConcern = concern
	}
}

// WithWriteTimeout sets the write timeout in milliseconds
func WithWriteTimeout(timeoutMS int) Option {
	return func(opts *MongoOptions) {
		opts.WTimeoutMS = timeoutMS
	}
}

// WithAuthSource sets the authentication database
// For root users created with MONGO_INITDB_ROOT_USERNAME, use "admin"
func WithAuthSource(authSource string) Option {
	return func(opts *MongoOptions) {
		opts.AuthSource = authSource
	}
}

// WithQueryLogging sets MongoDB query logging for debugging
func WithQueryLogging(enabled bool) Option {
	return func(opts *MongoOptions) {
		opts.QueryLogging = enabled
	}
}

// IsValidReadPreference checks if the read preference is valid
func IsValidReadPreference(pref string) bool {
	for _, valid := range validReadPreferences {
		if pref == valid {
			return true
		}
	}
	return false
}

// IsValidReadConcern checks if the read concern is valid
func IsValidReadConcern(concern string) bool {
	for _, valid := range validReadConcerns {
		if concern == valid {
			return true
		}
	}
	return false
}

// IsValidWriteConcern checks if the write concern is valid
func IsValidWriteConcern(concern string) bool {
	for _, valid := range validWriteConcerns {
		if concern == valid {
			return true
		}
	}
	return false
}
