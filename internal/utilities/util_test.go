package utilities_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bogdanutanu/go-rest-api-example/internal/models/data"
	"github.com/bogdanutanu/go-rest-api-example/internal/utilities"
	"github.com/stretchr/testify/assert"
)

func TestFormatTimeToISO(t *testing.T) {
	got := utilities.FormatTimeToISO(time.Date(2022, 5, 18, 9, 36, 0, 0, time.UTC))
	want := "2022-05-18T09:36:00Z"

	if want != got {
		t.Errorf("Expected '%s', but got '%s'", want, got)
	}
}

func TestCurrentISOTime(t *testing.T) {
	got := utilities.CurrentISOTime()
	parsedTime, err := time.Parse(time.RFC3339, got)
	if err != nil {
		t.Error("Received time string is not good format")
		return
	}
	z, offset := parsedTime.Zone()
	assert.Equal(t, "UTC", z)
	assert.Equal(t, 0, offset)
}

type DevModeTestCase struct {
	input  string
	result bool
}

func TestIsDevMode(t *testing.T) {
	cases := []DevModeTestCase{
		{"dev", true},
		{"development ", true},
		{"test", false},
		{"stage", false},
		{"production", false},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s=%t", tc.input, tc.result), func(t *testing.T) {
			got := utilities.IsDevMode(tc.input)
			if tc.result != got {
				t.Errorf("Expected '%t', but got '%t'", tc.result, got)
			}
		})
	}
}

func TestRandomPrice(t *testing.T) {
	price := utilities.RandomPrice()
	if price > utilities.MaxPrice {
		t.Errorf("Price is out of range: %v", price)
	}
}

func TestCalculateTotalAmount(t *testing.T) {
	// Test case: Empty input
	emptyTotal := utilities.CalculateTotalAmount([]data.Product{})
	// ugly hack to avoid lint error to use InEpsilon and InEpsilon limitation for 0 values
	assert.InEpsilon(t, 1, emptyTotal+1, 0)

	// Test case: Single product
	singleProductTotal := utilities.CalculateTotalAmount([]data.Product{
		{Name: "Product 1", Price: 10.0, Quantity: 1},
	})
	assert.InEpsilon(t, 10.0, singleProductTotal, 0.0001)

	// Test case: Multiple products
	multipleProductsTotal := utilities.CalculateTotalAmount([]data.Product{
		{Name: "Product 1", Price: 10.0, Quantity: 2},
		{Name: "Product 2", Price: 20.0, Quantity: 1},
		{Name: "Product 3", Price: 15.0, Quantity: 3},
	})
	assert.InEpsilon(t, 10.0*2+20.0*1+15.0*3, multipleProductsTotal, 0.0001)

	// Test case: Products with zero quantity
	zeroQuantityProductsTotal := utilities.CalculateTotalAmount([]data.Product{
		{Name: "Product 1", Price: 10.0, Quantity: 0},
		{Name: "Product 2", Price: 20.0, Quantity: 0},
	})
	// ugly hack to avoid lint error to use InEpsilon and InEpsilon limitation for 0 values
	assert.InEpsilon(t, 1, zeroQuantityProductsTotal+1, 0)
}
