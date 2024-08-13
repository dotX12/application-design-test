package config

import (
	"applicationDesignTest/config/api"
	"applicationDesignTest/config/otlp"
	_ "github.com/joho/godotenv/autoload"
)

type Factory struct {
	OTLP otlp.Config
	API  api.Config
}

func New() *Factory {
	cfg := Load()
	return cfg
}
