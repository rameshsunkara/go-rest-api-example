package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

const (
	DefaultMongoPort = 27017
	// MongoScheme is the standard MongoDB connection scheme
	MongoScheme = "mongodb"
	// MongoSRVScheme is the MongoDB SRV connection scheme used generally for cloud DBs
	MongoSRVScheme = "mongodb+srv"
	// DefaultMongoDBSidecar is the default path to the MongoDB sidecar file
	DefaultMongoDBSidecar string = "/secrets/db.json"
	// DefaultWriteTimeoutMS is the default write timeout in milliseconds for replica set operations
	DefaultWriteTimeoutMS = 5000
)

var (
	ErrNoHosts             = errors.New("at least one host is required")
	ErrSRVRequiresOneHost  = errors.New("SRV connection requires exactly one host")
	ErrInvalidPort         = errors.New("port must be between 1 and 65535")
	ErrInvalidReadPref     = errors.New("invalid read preference")
	ErrInvalidReadConcern  = errors.New("invalid read concern")
	ErrInvalidWriteConcern = errors.New("invalid write concern")
	ErrSideCarFileRead     = errors.New("failed to read sidecar file")
	ErrSideCarFileFormat   = errors.New("invalid sidecar file format")
)

// ConnectionURL creates a MongoDB connection URI using the given hosts and optional settings.
// The database parameter is optional - pass empty string if no database should be included in the URI.
// Example usage:
//
//	creds := &MongoCredentials{Username: "user", Password: "pass"}
//	uri, err := ConnectionURL([]string{"localhost"}, "mydb", creds,
//		WithPort(27017),
//		WithReplicaSet("rs0"),
//		WithReadPreference("secondary"))
func ConnectionURL(hosts []string, database string, creds *MongoCredentials, opts ...Option) (string, *MongoOptions, error) {
	if len(hosts) == 0 {
		return "", nil, ErrNoHosts
	}

	// Apply options
	options := applyOptions(opts...)

	if err := validateOptions(options, len(hosts)); err != nil {
		return "", nil, err
	}

	u := &url.URL{}

	// Set scheme
	u.Scheme = getScheme(options.UseSRV)

	// Set credentials
	if creds != nil && creds.Username != "" {
		u.User = getUserInfo(creds.Username, creds.Password)
	}

	// Set host
	hostString, err := getHost(hosts, options)
	if err != nil {
		return "", nil, err
	}
	u.Host = hostString

	// Set database path
	u.Path = getDatabasePath(database)

	// Set query parameters
	u.RawQuery = getQueryString(options)

	return u.String(), options, nil
}

// MaskConnectionURL takes an existing MongoDB connection URL and returns a version with masked credentials.
// Example: "mongodb://user:pass@host:27017/db" -> "mongodb://***:***@host:27017/db"
func MaskConnectionURL(connectionURL string) string {
	if connectionURL == "" {
		return ""
	}

	u, err := url.Parse(connectionURL)
	if err != nil {
		return connectionURL // Return original if parsing fails
	}

	// Mask the user info if it exists
	if u.User != nil {
		u.User = url.UserPassword("***", "***")
	}

	return u.String()
}

func MongoDBCredentialFromSideCar(sideCarFile string) (*MongoCredentials, error) {
	if sideCarFile == "" {
		sideCarFile = DefaultMongoDBSidecar
	}
	jsonFile, err := os.Open(sideCarFile)
	if err != nil {
		return nil, ErrSideCarFileRead
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, ErrSideCarFileRead
	}
	var mongoCredential MongoCredentials
	err = json.Unmarshal(byteValue, &mongoCredential)
	if err != nil {
		return nil, ErrSideCarFileFormat
	}
	return &mongoCredential, nil
}

// applyOptions applies the functional options to create MongoOptions with sensible defaults
func applyOptions(opts ...Option) *MongoOptions {
	// Start with sensible defaults
	options := &MongoOptions{
		Port:         DefaultMongoPort, // 27017
		UseSRV:       false,            // Standard connection by default
		QueryLogging: false,            // Disabled by default for performance
	}

	// Apply user-provided options to override defaults
	for _, opt := range opts {
		opt(options)
	}

	// If a replica set is configured, apply replica set specific defaults
	if options.ReplicaSet != "" {
		// Only set these defaults if not already specified by user options
		if options.ReadPreference == "" {
			options.ReadPreference = "primary"
		}
		if options.ReadConcern == "" {
			options.ReadConcern = "local"
		}
		if options.WriteConcern == "" {
			options.WriteConcern = "majority"
		}
		if options.WTimeoutMS == 0 {
			options.WTimeoutMS = DefaultWriteTimeoutMS // 5 second write timeout
		}
	}

	return options
}

// validateOptions validates the optional settings
func validateOptions(opts *MongoOptions, hostsCount int) error {
	if opts.UseSRV && hostsCount != 1 {
		return ErrSRVRequiresOneHost
	}

	if opts.Port != 0 && (opts.Port < 1 || opts.Port > 65535) {
		return ErrInvalidPort
	}

	if opts.ReadPreference != "" && !IsValidReadPreference(opts.ReadPreference) {
		return fmt.Errorf("%w: %s", ErrInvalidReadPref, opts.ReadPreference)
	}

	if opts.ReadConcern != "" && !IsValidReadConcern(opts.ReadConcern) {
		return fmt.Errorf("%w: %s", ErrInvalidReadConcern, opts.ReadConcern)
	}

	if opts.WriteConcern != "" && !IsValidWriteConcern(opts.WriteConcern) {
		return fmt.Errorf("%w: %s", ErrInvalidWriteConcern, opts.WriteConcern)
	}

	return nil
}

// getScheme returns the appropriate URI scheme based on UseSRV
func getScheme(useSRV bool) string {
	if useSRV {
		return MongoSRVScheme
	}
	return MongoScheme
}

// getUserInfo returns user credentials for URL if provided
func getUserInfo(username, password string) *url.Userinfo {
	if password != "" {
		return url.UserPassword(url.QueryEscape(username), url.QueryEscape(password))
	}
	return url.User(url.QueryEscape(username))
}

// getDatabasePath returns the database path for the URI
func getDatabasePath(database string) string {
	if database != "" {
		return "/" + database
	}
	return ""
}

// getHost returns the host portion of the MongoDB URI
func getHost(hosts []string, opts *MongoOptions) (string, error) {
	if opts.UseSRV {
		// SRV uses single host, no port (validation already done)
		return hosts[0], nil
	}

	// Standard connection - can have multiple hosts
	hostList := make([]string, 0, len(hosts))
	for _, h := range hosts {
		if strings.Contains(h, ":") {
			// Host already includes port
			hostList = append(hostList, h)
		} else if opts.Port > 0 {
			// Use configured port
			hostList = append(hostList, fmt.Sprintf("%s:%d", h, opts.Port))
		} else {
			// Use default MongoDB port
			hostList = append(hostList, fmt.Sprintf("%s:%d", h, DefaultMongoPort))
		}
	}

	return strings.Join(hostList, ","), nil
}

// getQueryString returns the query string for MongoDB URI
func getQueryString(opts *MongoOptions) string {
	q := url.Values{}

	if opts.ReplicaSet != "" {
		q.Set("replicaSet", opts.ReplicaSet)
	}
	if opts.ReadPreference != "" {
		q.Set("readPreference", opts.ReadPreference)
	}
	if opts.ReadConcern != "" {
		q.Set("readConcernLevel", opts.ReadConcern)
	}
	if opts.WriteConcern != "" {
		q.Set("w", opts.WriteConcern)
	}
	if opts.WTimeoutMS > 0 {
		q.Set("wtimeoutMS", fmt.Sprintf("%d", opts.WTimeoutMS))
	}
	if opts.AuthSource != "" {
		q.Set("authSource", opts.AuthSource)
	}

	return q.Encode()
}
