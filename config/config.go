package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerConfig  *ServerConfig
	StorageConfig *StorageConfig
	SecretConfig  *SecretConfig
}

type ServerConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

type StorageConfig struct {
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

type SecretConfig struct {
	UserKey string `env:"USER_KEY" envDefault:"_rj45"`
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

func NewSecretConfig() (*SecretConfig, error) {
	cfg := SecretConfig{}

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

	secretConfig, err := NewSecretConfig()

	if err != nil {
		return nil, err
	}

	return &Config{
		ServerConfig:  serverCfg,
		StorageConfig: storageCfg,
		SecretConfig:  secretConfig,
	}, nil
}

func (c *Config) ParseFlags() {
	a := flag.String("a", ":8080", "server address")
	b := flag.String("b", "http://localhost:8080", "base url")
	f := flag.String("f", "storage/filestorage/url_storage.json", "file path to storage")
	d := flag.String("d", "postgres://ya:url@localhost:5432", "database connection")

	flag.Parse()

	if c.ServerConfig.ServerAddress == "" {
		c.ServerConfig.ServerAddress = *a
	}

	if c.ServerConfig.BaseURL == "" {
		c.ServerConfig.BaseURL = *b
	}

	if c.StorageConfig.FileStoragePath == "" {
		c.StorageConfig.FileStoragePath = *f
	}

	if c.StorageConfig.DatabaseDSN == "" {
		c.StorageConfig.DatabaseDSN = *d
	}
}
