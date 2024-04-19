package util

import (
	"crypto/rand"
	"math/big"
	"strings"
	"time"
)

func FormatTimeToISO(timeToFormat time.Time) string {
	return timeToFormat.Format(time.RFC3339)
}

func CurrentISOTime() string {
	return FormatTimeToISO(time.Now().UTC())
}

// IsDevMode - Checks if the given string denotes any of the development environment.
func IsDevMode(s string) bool {
	return strings.Contains(s, "local") || strings.Contains(s, "dev")
}

// RandomPrice - Generates a random price between 0 and 1000.
const defaultPrice = 100
const maxPrice = 1000

func RandomPrice() uint64 {
	var price *big.Int
	var err error
	if price, err = rand.Int(rand.Reader, big.NewInt(maxPrice)); err != nil {
		price = big.NewInt(defaultPrice)
	}
	return price.Uint64()
}
