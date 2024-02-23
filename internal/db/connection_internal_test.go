package db

import (
	"errors"
	"testing"
)

func TestNewConnection(t *testing.T) {
	type newConnectionTestCase struct {
		Description string
		Input       string
		ExpectedErr error
	}

	var testCases = []newConnectionTestCase{
		{
			Description: "expect error when connection url is empty",
			Input:       "",
			ExpectedErr: ErrInvalidConnUrl,
		},
		{
			Description: "expect client creation error when connection url is invalid",
			Input:       "mongodb+srv://fuzzy-yogi:howsecureisthis@mongodb.net/?retryWrites=true&w=majority",
			ExpectedErr: ErrClientCreation,
		},
		{
			Description: "expect client object",
			Input:       "mongodb://test",
			ExpectedErr: ErrClientInit,
		},
	}

	for i, tc := range testCases {
		_, err := newClient(tc.Input)
		if !errors.Is(err, tc.ExpectedErr) {
			t.Errorf("TestNewConnection test case %d:%s failed: expected %v; got %v", i, tc.Description, tc.ExpectedErr, err)
		}
	}
}
