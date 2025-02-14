package config

import "os"

type Config struct {
	Port             string
	PostgresHost     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresPort     string
}

func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "8083"),
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresDB:       getEnv("POSTGRES_DB", "addresses"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
