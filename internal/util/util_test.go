package util_test

import (
	"fmt"
	"net/http"
	"net/url"
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

func TestHasUnSupportedQueryParams(t *testing.T) {
	testCases := []struct {
		description     string
		queryParams     url.Values
		supportedParams map[string]bool
		expectedVal     bool
	}{
		{
			description:     "All parameters are supported",
			queryParams:     url.Values{"param1": []string{"value1"}, "param2": []string{"value2"}},
			supportedParams: map[string]bool{"param1": true, "param2": true},
			expectedVal:     false,
		},
		{
			description:     "Some parameters are not supported",
			queryParams:     url.Values{"param1": []string{"value1"}, "param3": []string{"value3"}},
			supportedParams: map[string]bool{"param1": true, "param2": true},
			expectedVal:     true,
		},
		{
			description:     "No parameters are supported",
			queryParams:     url.Values{"param1": []string{"value1"}, "param3": []string{"value3"}},
			supportedParams: map[string]bool{},
			expectedVal:     true,
		},
		{
			description:     "handle when nil is passed as supportedParams",
			queryParams:     url.Values{"param1": []string{"value1"}, "param3": []string{"value3"}},
			supportedParams: nil,
			expectedVal:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			req := &http.Request{URL: &url.URL{RawQuery: tc.queryParams.Encode()}}
			supported := util.HasUnSupportedQueryParams(req, tc.supportedParams)
			if supported != tc.expectedVal {
				t.Errorf("Expected %v, but got %v", tc.expectedVal, supported)
			}
		})
	}
}

