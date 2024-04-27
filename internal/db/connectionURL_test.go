package db_test

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/rameshsunkara/go-rest-api-example/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestMongoDBCredentialFromSideCar(t *testing.T) {
	type mongoCredentialTestCase struct {
		Description string
		Input       string
		ExpectedOut *db.MongoDBCredentials
		ExpectedErr error
	}
	var testCases = []mongoCredentialTestCase{
		{
			Description: "ensure MongoDBCredentialFromSideCar returns expected MongoDBCredential",
			Input:       "../mockData/mongoDB_test_credentials.json",
			ExpectedOut: &db.MongoDBCredentials{
				Hostname:   "test",
				Password:   "123456789",
				User:       "ecuser",
				ReplicaSet: "",
				Port:       "8888",
			},
		},
		{
			Description: "ensure MongoDBCredentialFromSideCar returns error when invalid file path is given",
			Input:       "../mockData/non-existent.json",
			ExpectedErr: db.ErrSideCarFileRead,
		},
		{
			Description: "expect nil when invalid json file is given",
			Input:       "../mockData/mongoDB_test_credentials_fail.json",
			ExpectedOut: nil,
			ExpectedErr: db.ErrSideCarFileFormat,
		},
	}

	for i, tc := range testCases {
		t.Setenv("MongoVaultSideCar", tc.Input)
		result, err := db.MongoDBCredentialFromSideCar(os.Getenv("MongoVaultSideCar"))
		if !assert.EqualValues(t, tc.ExpectedOut, result) || !errors.Is(err, tc.ExpectedErr) {
			t.Errorf("MongoDBCredentialFromSideCar: Test Case: %d:%s failed: expected %v, %v; got %v, %v",
				i+1, tc.Description, tc.ExpectedOut, tc.ExpectedErr, result, err)
		}
	}
}

func TestMongoConnectionUrl(t *testing.T) {
	type connectionURLTestCase struct {
		Description string
		Input       *db.MongoDBCredentials
		ExpectedOut string
	}
	var testCases = []connectionURLTestCase{
		{
			Description: "ensure Connection Url is empty",
			Input:       &db.MongoDBCredentials{},
			ExpectedOut: "",
		},
		{
			Description: "atlas connection url should have options retryWrites and w",
			Input: &db.MongoDBCredentials{
				Hostname: "mongodb.net",
				User:     "fuzzy-yogi",
				Password: "howsecureisthis",
			},
			ExpectedOut: "mongodb+srv://fuzzy-yogi:howsecureisthis@mongodb.net/?" +
				"readPreference=nearest&retryWrites=true&w=majority",
		},
		{
			Description: "connection url should have username and password in url",
			Input: func() *db.MongoDBCredentials {
				var m db.MongoDBCredentials
				data, _ := os.ReadFile("../mockData/mongoDB_test_credentials.json")
				_ = json.Unmarshal(data, &m)
				return &m
			}(),
			ExpectedOut: "mongodb://ecuser:123456789@test:8888",
		},
		{
			Description: "should have replicaset",
			Input: &db.MongoDBCredentials{
				Hostname:   "mongodb1.svc.com,mongodb2.svc.com",
				ReplicaSet: "mySet",
			},
			ExpectedOut: "mongodb://mongodb1.svc.com,mongodb2.svc.com/?replicaSet=mySet",
		},
		{
			Description: "do not include authentication if password is missing",
			Input: &db.MongoDBCredentials{
				Hostname:   "mongodb1.svc.com,mongodb2.svc.com",
				ReplicaSet: "mySet",
				User:       "test",
			},
			ExpectedOut: "mongodb://mongodb1.svc.com,mongodb2.svc.com/?replicaSet=mySet",
		},
		{
			Description: "do not include authentication if username is missing",
			Input: &db.MongoDBCredentials{
				Hostname:   "mongodb1.svc.com,mongodb2.svc.com",
				ReplicaSet: "mySet",
				Password:   "test",
			},
			ExpectedOut: "mongodb://mongodb1.svc.com,mongodb2.svc.com/?replicaSet=mySet",
		},
		{
			Description: "include port when its given",
			Input: &db.MongoDBCredentials{
				Hostname:   "mongodb1.svc.com",
				ReplicaSet: "mySet",
				Password:   "test",
				Port:       "27107",
			},
			ExpectedOut: "mongodb://mongodb1.svc.com:27107/?replicaSet=mySet",
		},
		{
			Description: "discard port when its multiple hosts",
			Input: &db.MongoDBCredentials{
				Hostname:   "mongodb1.svc.com,mongodb2.svc.com",
				ReplicaSet: "mySet",
				Password:   "test",
				Port:       "23456",
			},
			ExpectedOut: "mongodb://mongodb1.svc.com,mongodb2.svc.com/?replicaSet=mySet",
		},
	}
	for i, tc := range testCases {
		result := db.MongoConnectionURL(tc.Input)
		if result != tc.ExpectedOut {
			t.Errorf("TestMongoConnectionUrl test case %d:%s failed: expected %s; got %s",
				i, tc.Description, tc.ExpectedOut, result)
		}
	}
}

func TestMaskedMongoConnectionUrl(t *testing.T) {
	type connectionURLTestCase struct {
		Description string
		Input       *db.MongoDBCredentials
		ExpectedOut string
	}
	var testCases = []connectionURLTestCase{
		{
			Description: "ensure masking doesnt fail for empty credentials",
			Input:       &db.MongoDBCredentials{},
			ExpectedOut: "",
		},
		{
			Description: "ensure username and passwords are masked",
			Input: &db.MongoDBCredentials{
				Hostname: "mongodb.net",
				User:     "fuzzy-yogi",
				Password: "howsecureisthis",
			},
			ExpectedOut: "mongodb+srv://#####:#####@mongodb.net/?readPreference=nearest&retryWrites=true&w=majority",
		},
	}

	for i, tc := range testCases {
		result := db.MaskedMongoConnectionURL(tc.Input)
		if result != tc.ExpectedOut {
			t.Errorf("TestMaskedMongoConnectionUrl test case %d:%s failed: expected %s; got %s",
				i, tc.Description, tc.ExpectedOut, result)
		}
	}
}
