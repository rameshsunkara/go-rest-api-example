package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	c, err := LoadConfig("dev")
	assert.NoError(t, err)
	k := c.AllKeys()
	assert.Equal(t, 2, len(k))
}

func TestLoadConfig_Failure(t *testing.T) {
	c, err := LoadConfig("dummy")
	assert.Error(t, err)
	assert.Equal(t, "", c.ConfigFileUsed())
}
