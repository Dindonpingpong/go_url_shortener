package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerConfig *ServerConfig
	StorageConfig *StorageConfig
}

type ServerConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"localhost:8080"`
}

type StorageConfig struct {
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"url_storage.json"`
}

func NewServerConfig() (*ServerConfig, error) {
	cfg := ServerConfig{}

	err := env.Parse(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func NewStorageConfig() (*StorageConfig, error) {
	cfg := StorageConfig{}

	err := env.Parse(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func NewDefaultConfiguration() (*Config, error) {
	serverCfg, err := NewServerConfig()

	if err != nil {
		return nil, err
	}

	storageCfg, err := NewStorageConfig()

	if err != nil {
		return nil, err
	}

	return &Config{
		ServerConfig: serverCfg,
		StorageConfig: storageCfg,
	}, nil
}

func (c *Config) ParseFlags() {
	a := flag.String("a", ":8080", "server address")
	b := flag.String("b", "http://localhost:8080", "base url")
	f := flag.String("f", "url_storage.json", "file path to storage")

	flag.Parse()

	c.ServerConfig.ServerAddress = *a
	c.ServerConfig.BaseURL = *b
	c.StorageConfig.FileStoragePath = *f
}