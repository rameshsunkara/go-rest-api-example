package flightrecorder

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime/trace"
	"time"

	"github.com/bogdanutanu/go-rest-api-example/pkg/logger"
)

// Default configuration values for the flight recorder.
const (
	// DefaultMinAge is the default minimum age of events to keep in the flight recorder buffer.
	DefaultMinAge = 1 * time.Second
	// DefaultMaxBytes is the default maximum size of the flight recorder buffer (1 MiB).
	DefaultMaxBytes uint64 = 1 << 20
	// DefaultTraceDir is the default directory where trace files are stored.
	DefaultTraceDir = "./traces"
)

// Recorder wraps the Go flight recorder with additional metadata and convenience methods.
// It maintains a rolling buffer of trace events that can be snapshotted on demand.
type Recorder struct {
	fr       *trace.FlightRecorder
	traceDir string
}

// New initializes and starts a flight recorder with the given configuration.
// Returns nil if initialization fails (e.g., directory creation error, flight recorder already enabled).
//
// Parameters:
//   - lgr: Logger for recording initialization events and errors
//   - traceDir: Directory where trace files will be stored (uses DefaultTraceDir if empty)
//   - minAge: Minimum age of events to keep in buffer (uses DefaultMinAge if 0)
//   - maxBytes: Maximum buffer size in bytes (uses DefaultMaxBytes if 0)
func New(lgr logger.Logger, traceDir string, minAge time.Duration, maxBytes uint64) *Recorder {
	// Use defaults if not provided
	if traceDir == "" {
		traceDir = DefaultTraceDir
	}
	if minAge == 0 {
		minAge = DefaultMinAge
	}
	if maxBytes == 0 {
		maxBytes = DefaultMaxBytes
	}

	// Ensure trace output directory exists
	if err := os.MkdirAll(traceDir, 0755); err != nil {
		lgr.Error().Err(err).Str("traceDir", traceDir).Msg("Failed to create trace output directory")
		return nil
	}

	// Set up the flight recorder
	fr := trace.NewFlightRecorder(trace.FlightRecorderConfig{
		MinAge:   minAge,
		MaxBytes: maxBytes,
	})
	if err := fr.Start(); err != nil {
		lgr.Error().Err(err).Msg("Failed to start flight recorder")
		return nil
	}

	lgr.Info().
		Dur("minAge", minAge).
		Interface("maxBytes", maxBytes).
		Str("traceDir", traceDir).
		Msg("Flight recorder initialized")

	return &Recorder{
		fr:       fr,
		traceDir: traceDir,
	}
}

// NewDefault initializes and starts a flight recorder with default configuration.
// This is a convenience wrapper around New() using DefaultTraceDir, DefaultMinAge, and DefaultMaxBytes.
func NewDefault(lgr logger.Logger) *Recorder {
	return New(lgr, DefaultTraceDir, DefaultMinAge, DefaultMaxBytes)
}

// TraceDir returns the directory where trace files are stored.
func (r *Recorder) TraceDir() string {
	return r.traceDir
}

// WriteTo writes a snapshot of the flight recorder's rolling buffer to the given writer.
// This is a low-level method that provides flexibility for custom trace handling.
//
// Use cases:
//   - Writing traces to custom destinations (HTTP responses, streams, etc.)
//   - Integration with external monitoring systems
//   - Custom trace processing pipelines
//
// For the common case of capturing slow requests to files, use CaptureSlowRequest instead.
func (r *Recorder) WriteTo(w io.Writer) (int64, error) {
	if r.fr == nil {
		return 0, nil
	}
	return r.fr.WriteTo(w)
}

// CaptureSlowRequest captures a flight recorder snapshot for a slow request and saves it to a file.
// This is a high-level convenience method that handles file creation, snapshot writing, and logging.
//
// The trace file is automatically named with the format: slow-request-{METHOD}-{PATH}-{TIMESTAMP}.trace
// and saved to the configured trace directory.
//
// Parameters:
//   - lgr: Logger for recording capture events and errors
//   - method: HTTP method of the slow request (e.g., "GET", "POST")
//   - path: URL path of the slow request (e.g., "/api/orders")
//   - elapsed: Duration the request took to complete
//
// Returns:
//   - The full path to the created trace file on success
//   - Empty string if capture fails (file creation error, write error, or nil recorder)
func (r *Recorder) CaptureSlowRequest(lgr logger.Logger, method, path string, elapsed time.Duration) string {
	if r.fr == nil {
		return ""
	}

	// Generate trace filename with request details
	traceFile := filepath.Join(
		r.traceDir,
		fmt.Sprintf("slow-request-%s-%s-%d.trace",
			method,
			extractPathSegment(path),
			time.Now().Unix(),
		),
	)

	// Create the trace file
	f, err := os.Create(traceFile)
	if err != nil {
		lgr.Error().
			Err(err).
			Str("traceFile", traceFile).
			Msg("Failed to create trace file for slow request")
		return ""
	}
	defer f.Close()

	// Write the flight recorder snapshot to the file
	if _, writeErr := r.WriteTo(f); writeErr != nil {
		lgr.Error().
			Err(writeErr).
			Str("traceFile", traceFile).
			Msg("Failed to write flight recorder snapshot")
		return ""
	}

	// Log successful capture
	lgr.Info().
		Str("method", method).
		Str("path", path).
		Dur("elapsed", elapsed).
		Str("traceFile", traceFile).
		Msg("Slow request detected - trace captured")

	return traceFile
}

// extractPathSegment extracts the last segment from a URL path for use in filenames.
// This ensures trace filenames are filesystem-safe and descriptive.
//
// Examples:
//   - "/api/orders" -> "orders"
//   - "/users/123" -> "123"
//   - "/" -> "root"
//   - "" -> "root"
func extractPathSegment(path string) string {
	const defaultSegment = "root"

	if path == "" {
		return defaultSegment
	}
	// Extract the last path segment
	segment := filepath.Base(path)
	if segment == "." || segment == "/" {
		return defaultSegment
	}
	// Only allow alphanumeric, hyphen, underscore in filename
	re := regexp.MustCompile("^[A-Za-z0-9-_]+$")
	if re.MatchString(segment) {
		return segment
	}
	return defaultSegment
}
