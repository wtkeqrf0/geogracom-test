package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerAddress      string `env:"SERVER_ADDRESS"`
	PostgresConnection string `env:"POSTGRES_CONNECTION"`
}

func GetConfig() *Config {
	var cfg Config

	flag.StringVar(&cfg.ServerAddress, "a", "0.0.0.0:8080", "HTTP server startup address")
	flag.StringVar(&cfg.PostgresConnection, "b", "postgres://postgres:1234@0.0.0.0:5432/geogracom", "Address of the database connection")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}
