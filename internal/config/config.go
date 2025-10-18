package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string

	SupabaseURL     string
	SupabaseAnonKey string

	JWTSecret string

	Port        string
	Environment string

	FrontendURL string
}

var AppConfig *Config

func Load() error {
	_ = godotenv.Load()

	config := &Config{
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		SupabaseURL:     getEnv("SUPABASE_URL", ""),
		SupabaseAnonKey: getEnv("SUPABASE_ANON_KEY", ""),
		JWTSecret:       getEnv("JWT_SECRET", ""),
		Port:            getEnv("PORT", "8080"),
		Environment:     getEnv("ENV", "development"),
		FrontendURL:     getEnv("FRONTEND_URL", "http://localhost:3000"),
	}

	if config.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if config.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if len(config.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	AppConfig = config
	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
