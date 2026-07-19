package handlers

import (
	"encoding/json"
	"testing"
)

func TestNormalizeFrontendActivityNavigation(t *testing.T) {
	activity, err := normalizeFrontendActivity(frontendActivityRequest{
		Type:      "navigation",
		FromPath:  "/portal?section=dashboard",
		ToPath:    "/portal?section=command-center",
		PageName:  "Command Center",
		FromPage:  "Dashboard",
		RouteName: "PortalDashboard",
		Section:   "command-center",
		Resource:  "Portal Navigation",
		Metadata: map[string]any{
			"token":         "must-not-log",
			"browser_title": "Command Center - NCS Uganda",
		},
	})
	if err != nil {
		t.Fatalf("normalizeFrontendActivity returned error: %v", err)
	}
	if activity.Action != "ui:navigate" {
		t.Fatalf("Action = %q, want ui:navigate", activity.Action)
	}
	if activity.EventType != "UI_NAVIGATION" {
		t.Fatalf("EventType = %q, want UI_NAVIGATION", activity.EventType)
	}
	if activity.Resource != "Portal Navigation" {
		t.Fatalf("Resource = %q, want Portal Navigation", activity.Resource)
	}
	if activity.ResourceID != "command-center" {
		t.Fatalf("ResourceID = %q, want command-center", activity.ResourceID)
	}
	if activity.Endpoint != "/portal?section=command-center" {
		t.Fatalf("Endpoint = %q, want dashboard path", activity.Endpoint)
	}
	var payload map[string]any
	if err := json.Unmarshal(activity.PayloadExcerpt, &payload); err != nil {
		t.Fatalf("payload unmarshal failed: %v", err)
	}
	if payload["from_path"] != "/portal?section=dashboard" || payload["to_path"] != "/portal?section=command-center" {
		t.Fatalf("payload paths = %#v", payload)
	}
	if _, ok := payload["token"]; ok {
		t.Fatalf("sensitive token metadata was retained: %#v", payload)
	}
}

func TestNormalizeFrontendActivityRejectsUnknownType(t *testing.T) {
	if _, err := normalizeFrontendActivity(frontendActivityRequest{Type: "password", ToPath: "/portal"}); err == nil {
		t.Fatal("expected unsupported frontend activity type to be rejected")
	}
}
