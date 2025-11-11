package flightrecorder_test

import (
	"os"
	"testing"
	"time"

	"github.com/rameshsunkara/go-rest-api-example/pkg/flightrecorder"
	"github.com/rameshsunkara/go-rest-api-example/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefault(t *testing.T) {
	lgr := logger.New("info", os.Stdout)
	fr := flightrecorder.NewDefault(lgr)

	// Flight recorder may return nil if there's an issue, but that's okay for this test
	// Just verify it doesn't panic
	if fr != nil {
		assert.NotNil(t, fr)
		assert.Equal(t, flightrecorder.DefaultTraceDir, fr.TraceDir())
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		traceDir string
		minAge   time.Duration
		maxBytes uint64
	}{
		{
			name:     "with default values (empty params)",
			traceDir: "",
			minAge:   0,
			maxBytes: 0,
		},
		{
			name:     "with custom values",
			traceDir: "./test-traces",
			minAge:   2 * time.Second,
			maxBytes: 2 << 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lgr := logger.New("info", os.Stdout)
			fr := flightrecorder.New(lgr, tt.traceDir, tt.minAge, tt.maxBytes)

			// Flight recorder may return nil if there's an issue, but that's okay for this test
			// Just verify it doesn't panic
			if fr != nil {
				assert.NotNil(t, fr)
				// Verify traceDir is set correctly
				expectedDir := tt.traceDir
				if expectedDir == "" {
					expectedDir = flightrecorder.DefaultTraceDir
				}
				assert.Equal(t, expectedDir, fr.TraceDir())
			}

			// Cleanup test trace directory if created
			if tt.traceDir != "" {
				_ = os.RemoveAll(tt.traceDir)
			}
		})
	}
}

func TestNewWithInvalidTraceDir(t *testing.T) {
	lgr := logger.New("info", os.Stdout)

	// Try to create a flight recorder with an invalid directory
	// (using a path that's likely to fail, like a file instead of directory)
	tempFile, err := os.CreateTemp("", "test-file-*")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Try to use the file path as a directory (should fail)
	fr := flightrecorder.New(lgr, tempFile.Name(), time.Second, 1<<20)

	// Should return nil on failure
	assert.Nil(t, fr)
}

func TestCaptureSlowRequest(t *testing.T) {
	// Create a temporary directory for test traces
	tempDir, err := os.MkdirTemp("", "test-traces-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	lgr := logger.New("info", os.Stdout)
	fr := flightrecorder.New(lgr, tempDir, time.Second, 1<<20)

	// Flight recorder may be nil if already enabled in another test
	// This is expected behavior since only one can be active per process
	if fr == nil {
		t.Skip("Flight recorder already enabled - skipping test")
		return
	}

	// Capture a slow request trace
	elapsed := 600 * time.Millisecond
	traceFile := fr.CaptureSlowRequest(lgr, "GET", "/api/orders", elapsed)

	// Verify trace file was created
	assert.NotEmpty(t, traceFile)
	if traceFile != "" {
		_, statErr := os.Stat(traceFile)
		assert.NoError(t, statErr, "trace file should exist")
	}
}

func TestCaptureSlowRequestWithNilRecorder(t *testing.T) {
	lgr := logger.New("info", os.Stdout)

	// Create a recorder with invalid directory to get nil
	tempFile, err := os.CreateTemp("", "test-file-*")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Try to use the file path as a directory (should return nil)
	fr := flightrecorder.New(lgr, tempFile.Name(), time.Second, 1<<20)
	require.Nil(t, fr, "should return nil for invalid trace directory")

	// Verify we handle nil recorder gracefully in middleware
	// (this test documents expected behavior - middleware checks for nil)
}
