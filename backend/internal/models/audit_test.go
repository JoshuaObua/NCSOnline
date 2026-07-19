package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestAuditLogMarshalAddsSpecificationAliases(t *testing.T) {
	userID := "user-1"
	createdAt := time.Date(2026, 7, 19, 1, 2, 3, 123000000, time.UTC)
	log := AuditLog{
		ID:            "event-1",
		UserID:        &userID,
		Action:        "auth:login",
		Resource:      "auth",
		ResourceID:    "session-target",
		IPAddress:     "197.239.5.14",
		UserAgent:     "Mozilla/5.0 Chrome/126.0",
		Method:        "POST",
		Endpoint:      "/api/v1/auth/login",
		ResponseCode:  200,
		EventStatus:   "SUCCESS",
		SeverityLevel: "INFO",
		GeoCountry:    "Uganda",
		GeoTimezone:   "Africa/Kampala",
		Platform:      "Windows 11",
		Browser:       "Chrome",
		OSName:        "Windows 11",
		ClientType:    "Browser",
		SessionID:     "correlation-1",
		Username:      "admin@ncs.go.ug",
		EntryHash:     "abc123",
		CreatedAt:     createdAt,
	}

	payload, err := json.Marshal(log)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(payload, &got); err != nil {
		t.Fatal(err)
	}

	want := map[string]string{
		"event_id":            "event-1",
		"correlation_id":      "correlation-1",
		"action":              "auth:login",
		"status":              "SUCCESS",
		"severity":            "INFO",
		"actor_id":            "user-1",
		"actor_type":          "USER",
		"actor_username":      "admin@ncs.go.ug",
		"client_ip":           "197.239.5.14",
		"access_medium":       "WEB_UI",
		"parsed_client_agent": "Browser",
		"parsed_os":           "Windows 11",
		"parsed_browser":      "Chrome",
		"request_url":         "/api/v1/auth/login",
		"http_method":         "POST",
		"resource_type":       "auth",
		"signature":           "abc123",
		"timezone":            "Africa/Kampala",
	}
	for key, value := range want {
		if got[key] != value {
			t.Fatalf("%s = %#v, want %q in %s", key, got[key], value, payload)
		}
	}
	if got["timestamp"] == "" {
		t.Fatalf("timestamp missing in %s", payload)
	}
}
