package util

import (
	"strings"
	"time"
)

func FormatTimeToISO(timeToFormat time.Time) string {
	return timeToFormat.Format(time.RFC3339)
}

func CurrentISOTime() string {
	return FormatTimeToISO(time.Now().UTC())
}

// IsDevMode - Checks if the given string denotes any of the development environment
func IsDevMode(s string) bool {
	return strings.Contains(s, "dev")
}
