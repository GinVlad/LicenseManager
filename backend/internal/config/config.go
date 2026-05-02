package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	AdminEmail  string
	AdminPass   string
	Env         string
	CORSOrigins []string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	corsOrigins := []string{}
	if raw := os.Getenv("CORS_ORIGINS"); raw != "" {
		for _, o := range strings.Split(raw, ",") {
			if s := strings.TrimSpace(o); s != "" {
				corsOrigins = append(corsOrigins, s)
			}
		}
	}

	dbUrl := mustEnv("DATABASE_URL")
	if !strings.Contains(dbUrl, "sslmode=") {
		if strings.Contains(dbUrl, "?") {
			dbUrl += "&sslmode=disable"
		} else {
			dbUrl += "?sslmode=disable"
		}
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: dbUrl,
		JWTSecret:   mustEnv("JWT_SECRET"),
		AdminEmail:  getEnv("ADMIN_EMAIL", "admin@example.com"),
		AdminPass:   getEnv("ADMIN_PASSWORD", ""),
		Env:         getEnv("ENV", "development"),
		CORSOrigins: corsOrigins,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %s is not set", key)
	}
	return v
}
