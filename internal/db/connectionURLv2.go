package db

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	DefaultMongoPort = 27017
	// MongoScheme is the standard MongoDB connection scheme
	MongoScheme = "mongodb"
	// MongoSRVScheme is the MongoDB SRV connection scheme used generally for cloud DBs
	MongoSRVScheme = "mongodb+srv"
)

var (
	ErrNoHosts             = errors.New("at least one host is required")
	ErrSRVRequiresOneHost  = errors.New("SRV connection requires exactly one host")
	ErrInvalidPort         = errors.New("port must be between 1 and 65535")
	ErrInvalidReadPref     = errors.New("invalid read preference")
	ErrInvalidWriteConcern = errors.New("invalid write concern")
)

// ConnectionURL creates a MongoDB connection URI using the given hosts and optional settings.
// Example usage:
//
//	uri, err := ConnectionURL([]string{"localhost"},
//		WithCredentials("user", "pass"),
//		WithDatabase("mydb"),
//		WithPort(27017),
//		WithReplicaSet("rs0"),
//		WithReadPreference("secondary"))
func ConnectionURL(hosts []string, opts ...Option) (string, error) {
	if len(hosts) == 0 {
		return "", ErrNoHosts
	}

	// Apply options
	options := applyOptions(opts...)

	if err := validateOptions(options, len(hosts)); err != nil {
		return "", err
	}

	u := &url.URL{}

	// Set scheme
	u.Scheme = getScheme(options.UseSRV)

	// Set credentials
	if options.Username != "" {
		u.User = getUserInfo(options.Username, options.Password)
	}

	// Set host
	hostString, err := getHost(hosts, options)
	if err != nil {
		return "", err
	}
	u.Host = hostString

	// Set database path
	u.Path = getDatabasePath(options.Database)

	// Set query parameters
	u.RawQuery = getQueryString(options)

	return u.String(), nil
}

// MaskConnectionURL takes an existing MongoDB connection URL and returns a version with masked credentials.
// This is useful for safely logging connection URLs that have already been built.
//
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

// applyOptions applies the functional options to create MongoOptions
func applyOptions(opts ...Option) *MongoOptions {
	options := &MongoOptions{}
	for _, opt := range opts {
		opt(options)
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
		return fmt.Errorf("invalid read concern: %s", opts.ReadConcern)
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

	return q.Encode()
}
