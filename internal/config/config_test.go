package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	err := LoadConfig("dev")
	assert.NoError(t, err)
	k := config.AllKeys()
	assert.Equal(t, 2, len(k))
}

func TestLoadConfig_Failure(t *testing.T) {
	err := LoadConfig("dummy")
	assert.Error(t, err)
}
