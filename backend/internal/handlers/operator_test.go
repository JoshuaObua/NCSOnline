package handlers

import "testing"

func TestNormalizeServiceStatusUsesActualOperationalStates(t *testing.T) {
	tests := []struct {
		name   string
		state  string
		health string
		want   string
	}{
		{name: "running healthy", state: "running", health: "healthy", want: "running"},
		{name: "running without health", state: "running", health: "", want: "running"},
		{name: "exited", state: "exited", health: "", want: "stopped"},
		{name: "restarting", state: "restarting", health: "", want: "restarting"},
		{name: "missing state", state: "", health: "unknown", want: "unavailable"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := normalizeServiceState(tc.state, tc.health); got != tc.want {
				t.Fatalf("normalizeServiceState(%q, %q) = %q, want %q", tc.state, tc.health, got, tc.want)
			}
		})
	}
}

func TestNormalizeServiceHealthNeverReturnsUnknown(t *testing.T) {
	tests := []struct {
		name   string
		state  string
		health string
		want   string
	}{
		{name: "stopped service", state: "exited", health: "unknown", want: "stopped"},
		{name: "running without health", state: "running", health: "", want: "no healthcheck"},
		{name: "running healthy", state: "running", health: "healthy", want: "healthy"},
		{name: "running unhealthy", state: "running", health: "unhealthy", want: "unhealthy"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := normalizeServiceHealth(tc.state, tc.health); got != tc.want {
				t.Fatalf("normalizeServiceHealth(%q, %q) = %q, want %q", tc.state, tc.health, got, tc.want)
			}
		})
	}
}

func TestOpsServiceCatalogUsesComposeServiceNames(t *testing.T) {
	for _, service := range []string{"nginx", "frontend", "backend", "postgres", "worker", "location-service"} {
		if !allowedOpsService(service) {
			t.Fatalf("%s should be allowed", service)
		}
	}
	for _, service := range []string{"nsmis-worker", "backup", "unknown"} {
		if allowedOpsService(service) {
			t.Fatalf("%s should not be allowed", service)
		}
	}
}
