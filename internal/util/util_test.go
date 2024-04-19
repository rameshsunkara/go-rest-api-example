package util_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rameshsunkara/go-rest-api-example/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestFormatTimeToISO(t *testing.T) {
	got := util.FormatTimeToISO(time.Date(2022, 5, 18, 9, 36, 0, 0, time.UTC))
	want := "2022-05-18T09:36:00Z"

	if want != got {
		t.Errorf("Expected '%s', but got '%s'", want, got)
	}
}

func TestCurrentISOTime(t *testing.T) {
	got := util.CurrentISOTime()
	parsedTime, err := time.Parse(time.RFC3339, got)
	z, offset := parsedTime.Zone()

	if err != nil {
		t.Error("Recieved time string is not good format")
	}
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
			got := util.IsDevMode(tc.input)
			if tc.result != got {
				t.Errorf("Expected '%t', but got '%t'", tc.result, got)
			}
		})
	}
}

func TestRandomPrice(t *testing.T) {
	price := util.RandomPrice()
	if price > util.MaxPrice {
		t.Errorf("Price is out of range: %v", price)
	}
}
