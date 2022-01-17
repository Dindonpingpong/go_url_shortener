package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerConfig *ServerConfig
}

type ServerConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"localhost:8080"`
}

func NewDefaultConfiguration() (*Config, error) {
	cfg := ServerConfig{}

	err := env.Parse(&cfg)

	if err != nil {
		return nil, err
	}

	return &Config{
		ServerConfig: &cfg,
	}, nil
}
