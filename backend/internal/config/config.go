package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv          string
	Port            string
	DatabaseURL     string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	AllowedOrigins  []string
	LogLevel        string
	SupabaseURL     string
	SupabaseKey     string
	StorageBucket   string
	RateLimitReqs   int
	RateLimitWindow time.Duration
}

func Load() (*Config, error) {
	dbURL, err := mustEnv("DATABASE_URL")
	if err != nil {
		return nil, err
	}
	jwtSecret, err := mustEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}

	rlReqs, _ := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "100"))
	if rlReqs <= 0 {
		rlReqs = 100
	}

	return &Config{
		AppEnv:          getEnv("APP_ENV", "development"),
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     dbURL,
		JWTSecret:       jwtSecret,
		AccessTokenTTL:  parseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m")),
		RefreshTokenTTL: parseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h")),
		AllowedOrigins:  strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		SupabaseURL:     getEnv("SUPABASE_URL", ""),
		SupabaseKey:     getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
		StorageBucket:   getEnv("SUPABASE_STORAGE_BUCKET", "ncsms-uploads"),
		RateLimitReqs:   rlReqs,
		RateLimitWindow: parseDuration(getEnv("RATE_LIMIT_WINDOW", "1m")),
	}, nil
}

func (c *Config) IsDevelopment() bool { return c.AppEnv == "development" }
func (c *Config) IsProduction() bool  { return c.AppEnv == "production" }

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("required environment variable not set: %s", key)
	}
	return v, nil
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 15 * time.Minute
	}
	return d
}
