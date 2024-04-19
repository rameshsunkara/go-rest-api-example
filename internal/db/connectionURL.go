package db

import (
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"os"
	"strings"
)

type MongoDBCredentials struct {
	User       string `json:"user"`
	Port       string `json:"port"`
	Hostname   string `json:"hostname"`
	ReplicaSet string `json:"replicaset"`
	Password   string `json:"password"`
}

var (
	ErrSideCarFileFormat = errors.New("mongo side car file is in invalid format")
	ErrSideCarFileRead   = errors.New("unable to read the mongo side car file")
)

const DefaultMongoDBSidecar string = "/vault/secrets/mongodb.json"
const mongoAtlasIdentifier string = "mongodb.net"

func MongoConnectionURL(mc *MongoDBCredentials) string {
	if len(mc.Hostname) == 0 {
		return ""
	}
	var protocol, authMechanism, host string
	connParams := url.Values{}
	protocol = "mongodb"

	// Atlas uses different options
	if strings.Contains(mc.Hostname, mongoAtlasIdentifier) {
		protocol = "mongodb+srv"
		connParams.Add("retryWrites", "true")
		connParams.Add("w", "majority")
		connParams.Add("readPreference", "nearest")
	}

	if len(mc.ReplicaSet) > 0 {
		connParams.Add("replicaSet", mc.ReplicaSet)
	}

	if len(mc.User) > 0 && len(mc.Password) > 0 {
		authMechanism = mc.User + ":" + mc.Password + "@"
	}

	host = mc.Hostname
	isMultiHost := len(strings.Split(mc.Hostname, ",")) > 1
	if len(mc.Port) > 0 && !isMultiHost {
		host += ":" + mc.Port
	}
	finalURL := protocol + "://" + authMechanism + host

	connParamsStr := connParams.Encode()
	if len(connParamsStr) > 0 {
		finalURL += "/?" + connParamsStr
	}

	return finalURL
}

func MaskedMongoConnectionURL(mc *MongoDBCredentials) string {
	if len(mc.User) > 0 {
		mc.User = "#####"
	}

	if len(mc.Password) > 0 {
		mc.Password = "#####"
	}

	return MongoConnectionURL(mc)
}

func MongoDBCredentialFromSideCar(sideCarFile string) (*MongoDBCredentials, error) {
	if sideCarFile == "" {
		sideCarFile = DefaultMongoDBSidecar
	}
	jsonFile, err := os.Open(sideCarFile)
	if err != nil {
		return nil, ErrSideCarFileRead
	}
	byteValue, _ := io.ReadAll(jsonFile)
	var mongoDBCredential MongoDBCredentials
	err = json.Unmarshal(byteValue, &mongoDBCredential)
	if err != nil {
		return nil, ErrSideCarFileFormat
	}
	return &mongoDBCredential, nil
}
