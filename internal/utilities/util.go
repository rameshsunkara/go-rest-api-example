package utilities

import (
	"crypto/rand"
	"math/big"
	"strings"
	"time"

	"github.com/bogdanutanu/go-rest-api-example/internal/models/data"
)

// FormatTimeToISO returns the time in RFC3339 format.
func FormatTimeToISO(timeToFormat time.Time) string {
	return timeToFormat.Format(time.RFC3339)
}

// CurrentISOTime returns the current UTC time in RFC3339 format.
func CurrentISOTime() string {
	return FormatTimeToISO(time.Now().UTC())
}

// IsDevMode - Checks if the given string denotes any of the development environment.
func IsDevMode(s string) bool {
	return strings.Contains(s, "local") || strings.Contains(s, "dev")
}

const (
	defaultPrice = 100
	MaxPrice     = 1000
)

// RandomPrice - Generates a random price between 0 and 1000.
func RandomPrice() float64 {
	var price *big.Int
	var err error
	if price, err = rand.Int(rand.Reader, big.NewInt(MaxPrice)); err != nil {
		price = big.NewInt(defaultPrice)
	}
	pf, _ := price.Float64()
	return pf
}

// CalculateTotalAmount calculates the total amount of the order based on the prices of products.
func CalculateTotalAmount(products []data.Product) float64 {
	var total float64
	for _, product := range products {
		total += product.Price * (float64(product.Quantity))
	}
	return total
}
