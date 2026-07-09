package storage

import (
	"encoding/json"
	"strings"
)

const SettingsKey = "file_storage"

type Settings struct {
	PublicProvider      string `json:"public_provider"`
	ApplicationProvider string `json:"application_provider"`

	LocalPublicPath      string `json:"local_public_path"`
	LocalPublicURLPrefix string `json:"local_public_url_prefix"`
	LocalAppPath         string `json:"local_app_path"`
	LocalAppURLPrefix    string `json:"local_app_url_prefix"`

	GoogleDriveFolderID     string `json:"google_drive_folder_id"`
	GoogleDriveClientID     string `json:"google_drive_client_id,omitempty"`
	GoogleDriveClientSecret string `json:"google_drive_client_secret,omitempty"`
	GoogleDriveRefreshToken string `json:"google_drive_refresh_token,omitempty"`
	GoogleDriveMakePublic   bool   `json:"google_drive_make_public"`

	S3Bucket          string `json:"s3_bucket"`
	S3Region          string `json:"s3_region"`
	S3Prefix          string `json:"s3_prefix"`
	S3Endpoint        string `json:"s3_endpoint"`
	S3PublicBaseURL   string `json:"s3_public_base_url"`
	S3ForcePathStyle  bool   `json:"s3_force_path_style"`
	S3AccessKeyID     string `json:"s3_access_key_id,omitempty"`
	S3SecretAccessKey string `json:"s3_secret_access_key,omitempty"`

	SupabaseURL        string `json:"supabase_url"`
	SupabaseBucket     string `json:"supabase_bucket"`
	SupabasePrefix     string `json:"supabase_prefix"`
	SupabasePublicRead bool   `json:"supabase_public_read"`
	SupabaseServiceKey string `json:"supabase_service_key,omitempty"`
}

// allowedProviders is the single set of provider keys an admin may pick for
// EITHER scope (public portal assets or application evidence uploads), so
// switching providers works the same way everywhere file uploads happen.
var allowedProviders = map[string]bool{"local": true, "google_drive": true, "s3": true, "supabase": true}

func DefaultSettings() Settings {
	return Settings{
		PublicProvider:        "google_drive",
		ApplicationProvider:   "s3",
		LocalPublicPath:       "/app/uploads",
		LocalPublicURLPrefix:  "/uploads",
		LocalAppPath:          "/app/uploads/applications",
		LocalAppURLPrefix:     "/uploads/applications",
		GoogleDriveMakePublic: true,
		S3Region:              "us-east-1",
		S3Prefix:              "applications",
	}
}

func ParseSettings(raw []byte) (Settings, error) {
	cfg := DefaultSettings()
	if len(raw) == 0 {
		return cfg, nil
	}
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return cfg, err
	}
	cfg.Normalize()
	return cfg, nil
}

func (s *Settings) Normalize() {
	s.PublicProvider = normalizeProvider(s.PublicProvider, "google_drive", allowedProviders)
	s.ApplicationProvider = normalizeProvider(s.ApplicationProvider, "s3", allowedProviders)
	if strings.TrimSpace(s.LocalPublicPath) == "" {
		s.LocalPublicPath = "/app/uploads"
	}
	if strings.TrimSpace(s.LocalPublicURLPrefix) == "" {
		s.LocalPublicURLPrefix = "/uploads"
	}
	if strings.TrimSpace(s.LocalAppPath) == "" {
		s.LocalAppPath = "/app/uploads/applications"
	}
	if strings.TrimSpace(s.LocalAppURLPrefix) == "" {
		s.LocalAppURLPrefix = "/uploads/applications"
	}
	if strings.TrimSpace(s.S3Region) == "" {
		s.S3Region = "us-east-1"
	}
	if strings.TrimSpace(s.S3Prefix) == "" {
		s.S3Prefix = "applications"
	}
}

func normalizeProvider(value, fallback string, allowed map[string]bool) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if allowed[value] {
		return value
	}
	return fallback
}

func (s Settings) Redacted() Settings {
	s.GoogleDriveClientSecret = ""
	s.GoogleDriveRefreshToken = ""
	s.S3AccessKeyID = redact(s.S3AccessKeyID)
	s.S3SecretAccessKey = ""
	s.SupabaseServiceKey = ""
	return s
}

// SupabaseConfigured reports whether a service role key has been saved.
// The key itself is never sent back to the client (see Redacted), so the
// admin UI needs this boolean to know a key is already on file.
func (s Settings) SupabaseConfigured() bool {
	return strings.TrimSpace(s.SupabaseServiceKey) != ""
}

// GoogleDriveConnected reports whether a refresh token has been obtained via
// the OAuth connect flow, i.e. uploads can actually authenticate.
func (s Settings) GoogleDriveConnected() bool {
	return strings.TrimSpace(s.GoogleDriveRefreshToken) != ""
}

func (s Settings) MergeSecrets(existing Settings) Settings {
	if strings.TrimSpace(s.GoogleDriveClientSecret) == "" {
		s.GoogleDriveClientSecret = existing.GoogleDriveClientSecret
	}
	// The refresh token is only ever set internally by the OAuth connect
	// callback, never by the settings form, so preserve it here — unless the
	// client ID/secret changed, in which case the old token belongs to a
	// different OAuth app and can no longer be used; force a reconnect.
	if s.GoogleDriveClientID == existing.GoogleDriveClientID && s.GoogleDriveClientSecret == existing.GoogleDriveClientSecret {
		s.GoogleDriveRefreshToken = existing.GoogleDriveRefreshToken
	}
	if strings.TrimSpace(s.S3AccessKeyID) == "" || strings.HasPrefix(s.S3AccessKeyID, "********") {
		s.S3AccessKeyID = existing.S3AccessKeyID
	}
	if strings.TrimSpace(s.S3SecretAccessKey) == "" {
		s.S3SecretAccessKey = existing.S3SecretAccessKey
	}
	if strings.TrimSpace(s.SupabaseServiceKey) == "" {
		s.SupabaseServiceKey = existing.SupabaseServiceKey
	}
	return s
}

func redact(value string) string {
	if value == "" {
		return ""
	}
	if len(value) <= 4 {
		return "********"
	}
	return "********" + value[len(value)-4:]
}
