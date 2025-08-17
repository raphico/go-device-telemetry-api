package config

import (
	"os"
	"time"
)

type Config struct {
	DatabaseURL     string
	HTTPAddr        string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Env             string
}

func Load() Config {
	accessTokenTTL := 15 * time.Minute
	refreshTokenTTL := 7 * 24 * time.Hour

	if v := os.Getenv("ACCESS_TOKEN_TTL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			accessTokenTTL = d
		}
	}

	if v := os.Getenv("REFRESH_TOKEN_TTL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			refreshTokenTTL = d
		}
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/dbname"
	}

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable is required")
	}

	return Config{
		DatabaseURL:     dbURL,
		HTTPAddr:        ":" + port,
		JWTSecret:       jwtSecret,
		RefreshTokenTTL: refreshTokenTTL,
		AccessTokenTTL:  accessTokenTTL,
		Env:             env,
	}
}
