package util

import (
	"log/slog"
	"os"
)

type Config struct {
	Port     string
	BindAddr string
	LogLevel slog.Level
}

func LoadConfig() Config {
	return Config{
		Port:     GetEnv("PORT", "9000"),
		BindAddr: GetEnv("BIND_ADDR", ""),
		LogLevel: ParseLogLevel(GetEnv("LOG_LEVEL", "info")),
	}
}

func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
