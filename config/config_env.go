package config

import "github.com/caarlos0/env/v6"

type Config struct {
	Mongo MongoConfig
}

type MongoConfig struct {
	Host     string `env:"MONGO_HOST"`
	Port     int64  `env:"MONGO_PORT"`
	Database string `env:"MONGO_DATABASE"`
	Username string `env:"MONGO_USERNAME"`
	Password string `env:"MONGO_PWD"`
}

func NewFromEnv() (*Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
