package config

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
)

type Configuration struct {
	BaseURL     string `env:"BASE_URL"`
	ServerAddr  string `env:"SERVER_ADDR"`
	File        string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN string `env:"DATABASE_DSN"`
}

func NewConfig() *Configuration {
	config := new(Configuration)

	// command line flags with min priority
	flag.StringVar(&config.ServerAddr, "a", ":8080", "server address")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "base url for short links")
	flag.StringVar(&config.DatabaseDSN, "d", "", "database dsn")
	flag.StringVar(&config.File, "f", "", "file to store links")
	flag.Parse()

	// environment variables with max priority
	err := env.Parse(config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
