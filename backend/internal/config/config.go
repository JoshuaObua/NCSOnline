package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv               string
	Port                 string
	DatabaseURL          string
	JWTSecret            string
	AccessTokenTTL       time.Duration
	RefreshTokenTTL      time.Duration
	AllowedOrigins       []string
	LogLevel             string
	SupabaseURL          string
	SupabaseKey          string
	StorageBucket        string
	PrivateStoragePath   string
	MaxEvidenceBytes     int64
	RateLimitReqs        int
	RateLimitWindow      time.Duration
	USSDCallbackSecret   string
	LocationServiceURL   string
	LocationServiceToken string
	LocationTimeout      time.Duration
	GeoAllowedCountries  []string
	GeoFailClosed        bool
	SMTPHost             string
	SMTPPort             string
	SMTPUser             string
	SMTPPassword         string
	SMTPFrom             string
	PublicAppURL         string
	LogDir               string
	MaintenanceEnabled   bool
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
		AppEnv:               getEnv("APP_ENV", "development"),
		Port:                 getEnv("PORT", "8080"),
		DatabaseURL:          dbURL,
		JWTSecret:            jwtSecret,
		AccessTokenTTL:       parseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m")),
		RefreshTokenTTL:      parseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h")),
		AllowedOrigins:       strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
		SupabaseURL:          getEnv("SUPABASE_URL", ""),
		SupabaseKey:          getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
		StorageBucket:        getEnv("SUPABASE_STORAGE_BUCKET", "ncsms-uploads"),
		PrivateStoragePath:   getEnv("PRIVATE_STORAGE_PATH", "./private-data"),
		MaxEvidenceBytes:     int64Env("MAX_EVIDENCE_BYTES", 15<<20),
		RateLimitReqs:        rlReqs,
		RateLimitWindow:      parseDuration(getEnv("RATE_LIMIT_WINDOW", "1m")),
		USSDCallbackSecret:   getEnv("USSD_CALLBACK_SECRET", ""),
		LocationServiceURL:   getEnv("LOCATION_SERVICE_URL", "http://location-service:8090"),
		LocationServiceToken: getEnv("LOCATION_INTERNAL_TOKEN", ""),
		LocationTimeout:      parseDuration(getEnv("LOCATION_SERVICE_TIMEOUT", "3s")),
		GeoAllowedCountries:  splitUpper(getEnv("GEO_ALLOWED_COUNTRIES", "UG")),
		GeoFailClosed:        strings.EqualFold(getEnv("GEO_FAIL_CLOSED", "false"), "true"),
		SMTPHost:             getEnv("SMTP_HOST", ""),
		SMTPPort:             getEnv("SMTP_PORT", "587"),
		SMTPUser:             getEnv("SMTP_USER", ""),
		SMTPPassword:         getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:             getEnv("SMTP_FROM", "NCS <noreply@ncs.go.ug>"),
		PublicAppURL:         strings.TrimRight(getEnv("PUBLIC_APP_URL", "http://localhost:3000"), "/"),
		LogDir:               getEnv("LOG_DIR", "/var/log/app"),
		MaintenanceEnabled:   !strings.EqualFold(getEnv("MAINTENANCE_MIDDLEWARE_ENABLED", "true"), "false"),
	}, nil
}

func splitUpper(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.ToUpper(strings.TrimSpace(part)); part != "" {
			out = append(out, part)
		}
	}
	return out
}

func int64Env(key string, fallback int64) int64 {
	v, err := strconv.ParseInt(getEnv(key, ""), 10, 64)
	if err != nil || v <= 0 {
		return fallback
	}
	return v
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
