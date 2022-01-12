package config

type Config struct {
	Port string
}

func NewDefaultConfiguration() *Config {
	return &Config{
		Port: ":8080",
	}
}