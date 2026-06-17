package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	SupabaseJWTSecret  string
	AdminKey           string
	GooglePlacesAPIKey string
}

func Load() (Config, error) {
	_ = loadDotEnv()

	cfg := Config{
		Port:               envOr("PORT", "8080"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		SupabaseJWTSecret:  os.Getenv("SUPABASE_JWT_SECRET"),
		AdminKey:           os.Getenv("ADMIN_KEY"),
		GooglePlacesAPIKey: os.Getenv("GOOGLE_PLACES_API_KEY"),
	}

	if cfg.DatabaseURL == "" {
		return cfg, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func loadDotEnv() error {
	if _, err := os.Stat(".env"); err != nil {
		return nil
	}
	// godotenv loaded in main if present
	return nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func ParseIntQuery(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return n
}

func ParseFloatQuery(s string) (float64, bool) {
	if s == "" {
		return 0, false
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return f, true
}
