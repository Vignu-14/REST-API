package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv" 
)

type Config struct {
	Port        string
	DatabaseUrl string
	JWTSecret   string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		return nil, errors.New("database url is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("jwt secret is required")
	}

	return &Config{
		Port:        port,
		DatabaseUrl: dbUrl,
		JWTSecret:   jwtSecret,
	}, nil
}
