package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var config *viper.Viper

// LoadConfig TODO: Add more profiles and obey all 12-Factor app rules
func LoadConfig(env string) {
	log.Debug().Str("env", env).Msg("loading config for")
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")
	err = config.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error occurred while reading config")
	}
}

func GetConfig() *viper.Viper {
	return config
}
