package handlers

import "testing"

func TestNormalizeServiceStatusUsesRunningOrIdle(t *testing.T) {
	tests := []struct {
		name   string
		state  string
		health string
		want   string
	}{
		{name: "running healthy", state: "running", health: "healthy", want: "running"},
		{name: "running without health", state: "running", health: "", want: "running"},
		{name: "exited", state: "exited", health: "", want: "idle"},
		{name: "missing state", state: "", health: "unknown", want: "idle"},
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
		{name: "idle service", state: "exited", health: "unknown", want: "idle"},
		{name: "running without health", state: "running", health: "", want: "healthy"},
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
