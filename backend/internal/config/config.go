package config

import (
	"os"
	"strings"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	CorsOrigins   []string
	MigrationsDir string
}

func Load() Config {
	return Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		CorsOrigins:   splitCSV(getEnv("CORS_ORIGINS", "")),
		MigrationsDir: getEnv("MIGRATIONS_DIR", "./migrations"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func splitCSV(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
