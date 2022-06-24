package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var config *viper.Viper

// LoadConfig This is just to demonstrate usage of viper.
// TODO: Add more profiles and obey all 12-Factor app rules
func LoadConfig(env string) (*viper.Viper, error) {
	log.Debug().Str("env", env).Msg("loading config")
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("../config/")
	config.AddConfigPath("../../config/")
	config.AddConfigPath("config/")
	err := config.ReadInConfig()
	if err != nil {
		log.Error().Err(err).Msg("error occurred while reading config")
	}
	return config, err
}

func GetConfig() *viper.Viper {
	return config
}
