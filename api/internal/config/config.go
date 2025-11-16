package config

import "os"

type Config struct {
	DBDSN   string
	Port    string
	BaseURL string
}

func Load() *Config {
	return &Config{
		DBDSN:   os.Getenv("DB_DSN"),
		Port:    os.Getenv("PORT"),
		BaseURL: os.Getenv("BASE_URL"),
	}
}
