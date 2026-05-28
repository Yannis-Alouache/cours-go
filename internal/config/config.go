package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Server struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	TokenTTL    int
}

func LoadServer() (Server, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Server{}, err
	}

	cfg := Server{
		Port:        getenv("PORT", "8080"),
		DatabaseURL: getenv("DATABASE_URL", getenv("APP_DATABASE_URL", "")),
		JWTSecret:   getenv("JWT_SECRET", getenv("APP_JWT_SECRET", "")),
		TokenTTL:    24,
	}

	switch {
	case cfg.DatabaseURL == "":
		return Server{}, errors.New("missing DATABASE_URL (or APP_DATABASE_URL)")
	case cfg.JWTSecret == "":
		return Server{}, errors.New("missing JWT_SECRET (or APP_JWT_SECRET)")
	default:
		return cfg, nil
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
