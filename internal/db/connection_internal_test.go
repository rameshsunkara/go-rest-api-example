package db

import "testing"

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
			ExpectedErr: InvalidConnUrlErr,
		},
		{
			Description: "expect client creation error when connection url is invalid",
			Input:       "mongodb+srv://fuzzy-yogi:howsecureisthis@mongodb.net/?retryWrites=true&w=majority",
			ExpectedErr: ClientCreationErr,
		},
		{
			Description: "expect client object",
			Input:       "mongodb://test",
			ExpectedErr: nil,
		},
	}

	for i, tc := range testCases {
		_, err := newClient(tc.Input)
		if err != tc.ExpectedErr {
			t.Errorf("TestNewConnection test case %d:%s failed: expected %v; got %v", i, tc.Description, tc.ExpectedErr, err)
		}
	}
}
