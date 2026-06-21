package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// LocationResult is supplied by the internal Python service. Identity fields
// are correlation data only; the Go JWT middleware remains authoritative.
type LocationResult struct {
	IP            string   `json:"ip"`
	Country       string   `json:"country"`
	CountryCode   string   `json:"country_code"`
	Region        string   `json:"region"`
	City          string   `json:"city"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
	Timezone      string   `json:"timezone"`
	ISP           string   `json:"isp"`
	IsProxy       bool     `json:"is_proxy"`
	IsVPN         bool     `json:"is_vpn"`
	IsTor         bool     `json:"is_tor"`
	IsHosting     bool     `json:"is_hosting"`
	Platform      string   `json:"platform"`
	Browser       string   `json:"browser"`
	DeviceType    string   `json:"device_type"`
	Authenticated bool     `json:"authenticated"`
	UserID        *string  `json:"user_id"`
	Source        string   `json:"source"`
}

type LocationClient struct {
	baseURL string
	token   string
	client  *http.Client
}

func NewLocationClient(baseURL, token string, timeout time.Duration) *LocationClient {
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	return &LocationClient{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		token:   token,
		client:  &http.Client{Timeout: timeout},
	}
}

func (c *LocationClient) Locate(ctx context.Context, ip, userAgent, userID string) (*LocationResult, error) {
	if c == nil || c.baseURL == "" {
		return nil, fmt.Errorf("location service is not configured")
	}
	body, err := json.Marshal(map[string]any{
		"ip": ip, "user_agent": userAgent,
		"authenticated": userID != "", "user_id": nullableLocationUserID(userID),
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/locate", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("X-Internal-Token", c.token)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("location service returned %d", resp.StatusCode)
	}
	var result LocationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func nullableLocationUserID(userID string) any {
	if strings.TrimSpace(userID) == "" {
		return nil
	}
	return userID
}
