package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL       string
	ServerAddress     string
	JWTSecret         string
	EncryptionKey     string
	MaxStoragePerUser int64
}

func NewConfig() (*Config, error) {
	maxStorage, err := strconv.ParseInt(os.Getenv("MAX_STORAGE_PER_USER"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_STORAGE_PER_USER: %w", err)
	}

	return &Config{
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		ServerAddress:     os.Getenv("SERVER_ADDRESS"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		EncryptionKey:     os.Getenv("ENCRYPTION_KEY"),
		MaxStoragePerUser: maxStorage,
	}, nil
}
