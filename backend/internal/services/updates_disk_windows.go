//go:build windows

package services

// Disk telemetry is best-effort and only relevant on the deployment host
// (Linux container). The Windows stub keeps `go build` happy on dev
// machines without pulling in a winapi dependency.
func diskInfo() DiskInfo {
	return DiskInfo{Healthy: true}
}
