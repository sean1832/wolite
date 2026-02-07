package env

import (
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"
	"wolite/internal/auth"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret    string
	DatabasePath string
	JWTExpiry    time.Duration
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using default values")
	} else {
		slog.Info("Loaded .env file")
	}

	jwtToken := os.Getenv("JWT_SECRET")
	var err error
	if jwtToken == "" {
		// if user not provided a JWT secret, generate a random one
		jwtToken, err = auth.GenerateRandomString(32)
		if err != nil {
			log.Fatalf("failed to generate JWT secret: %v", err)
		}
		slog.Warn("JWT_SECRET not provided, generated a random one")
	}
	jwtExpiryString := os.Getenv("JWT_EXPIRY_SECONDS")
	if jwtExpiryString == "" {
		log.Fatalf("JWT_EXPIRY_SECONDS not provided")
	}
	jwtExpiry, err := strconv.Atoi(jwtExpiryString)
	if err != nil {
		log.Fatalf("failed to parse JWT_EXPIRY_SECONDS: %v", err)
	}
	if jwtExpiry < 0 {
		log.Fatalf("JWT_EXPIRY_SECONDS must be a positive integer")
	}
	slog.Info("JWT_EXPIRY_SECONDS", "value", strconv.Itoa(jwtExpiry))

	return &Config{
		JWTSecret:    jwtToken,
		DatabasePath: os.Getenv("DATABASE_PATH"),
		JWTExpiry:    time.Duration(jwtExpiry) * time.Second,
	}
}
