package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	Port         string `mapstructure:"SERVER_PORT"`
	AllowedHosts string `mapstructure:"ALLOWED_HOSTS"`
}

type DatabaseConfiguration struct {
	DatabaseName string `mapstructure:"DATABASE_NAME"`
}

type Config struct {
	Server   ServerConfiguration   `mapstructure:",squash"`
	Database DatabaseConfiguration `mapstructure:",squash"`
}

func Load() (*Config, error) {
	var result map[string]interface{}
	var config Config
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Error("error in reading from .env ", err)
		return nil, err
	}

	err := viper.Unmarshal(&result)
	if err != nil {
		logrus.Errorf("Unable to decode into map, %v", err)
		return nil, err
	}

	err = mapstructure.Decode(result, &config)
	if err != nil {
		logrus.Error("error in decoding result ", err)
		return nil, err
	}
	return &config, nil
}
