package jwtfx

import (
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	SecretKey string
	Issuer    string
	Expire    time.Duration
}

func NewConfig() *Config {
	return &Config{
		SecretKey: os.Getenv("JWT_SECRET"),
		Issuer:    os.Getenv("JWT_ISSUER"),
		Expire:    24 * time.Hour,
	}
}
