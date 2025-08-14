package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	HTTPAddr    string
}

func Load() Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/dbname"
	}

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		DatabaseURL: dbURL,
		HTTPAddr:    ":" + port,
	}
}
