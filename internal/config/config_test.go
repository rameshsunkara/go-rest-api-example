package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfig_Failure(t *testing.T) {
	LoadConfig("dummy")
	assert.Nil(t, config)
}

func TestLoadConfig_Success(t *testing.T) {
	LoadConfig("dev")
	assert.Nil(t, config)
}
