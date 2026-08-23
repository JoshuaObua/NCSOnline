package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type IoTecConfig struct {
	BaseURL      string
	AuthURL      string
	ClientID     string
	ClientSecret string
	WalletID     string
}

type IoTecTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type IoTecCollectionRequest struct {
	Category    string  `json:"category"`
	Currency    string  `json:"currency"`
	WalletID    string  `json:"walletId"`
	ExternalID  string  `json:"externalId,omitempty"`
	Payer       string  `json:"payer"`
	PayerName   string  `json:"payerName,omitempty"`
	PayerNote   string  `json:"payerNote,omitempty"`
	Amount      float64 `json:"amount"`
	PayeeNote   string  `json:"payeeNote,omitempty"`
	Channel     string  `json:"channel,omitempty"`
	RedirectURL string  `json:"redirectUrl,omitempty"`
}

type IoTecCollectionResponse struct {
	ID             string  `json:"id"`
	CreatedAt      string  `json:"createdAt"`
	Category       string  `json:"category"`
	Status         string  `json:"status"`
	PaymentChannel string  `json:"paymentChannel"`
	StatusCode     string  `json:"statusCode"`
	StatusMessage  string  `json:"statusMessage"`
	ExternalID     string  `json:"externalId"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
}

type IoTecService struct {
	cfg        IoTecConfig
	httpClient *http.Client

	mu        sync.RWMutex
	token     string
	expiresAt time.Time
}

func NewIoTecService() *IoTecService {
	baseURL := strings.TrimRight(getEnvFallback("IOTEC_BASE_URL", "https://pay.iotec.io"), "/")
	authURL := strings.TrimRight(getEnvFallback("IOTEC_AUTH_URL", "https://id.iotec.io"), "/")
	clientID := getEnvFallback("IOTEC_CLIENT_ID", "xmen-client")
	clientSecret := getEnvFallback("IOTEC_CLIENT_SECRET", "LvEDsfLrKQjJ6fJsn3iwL71Lzf82ibVL3f6Jef5D")
	walletID := getEnvFallback("IOTEC_WALLET_ID", "c9f41d7b-3784-42ca-9cfe-76f4244e9ed3")

	return &IoTecService{
		cfg: IoTecConfig{
			BaseURL:      baseURL,
			AuthURL:      authURL,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			WalletID:     walletID,
		},
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func getEnvFallback(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// GetAccessToken fetches or returns cached OAuth2 Bearer token from id.iotec.io
func (s *IoTecService) GetAccessToken(ctx context.Context) (string, error) {
	s.mu.RLock()
	if s.token != "" && time.Now().Before(s.expiresAt.Add(-30*time.Second)) {
		token := s.token
		s.mu.RUnlock()
		return token, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Double check after lock
	if s.token != "" && time.Now().Before(s.expiresAt.Add(-30*time.Second)) {
		return s.token, nil
	}

	// If no client credentials configured, support sandbox simulation mode
	if s.cfg.ClientID == "" || s.cfg.ClientSecret == "" {
		s.token = "sandbox-token-" + time.Now().Format("20060102150405")
		s.expiresAt = time.Now().Add(1 * time.Hour)
		return s.token, nil
	}

	data := url.Values{}
	data.Set("client_id", s.cfg.ClientID)
	data.Set("client_secret", s.cfg.ClientSecret)
	data.Set("grant_type", "client_credentials")

	tokenURL := fmt.Sprintf("%s/connect/token", s.cfg.AuthURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("iotec auth network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("iotec auth failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp IoTecTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode iotec token response: %w", err)
	}

	s.token = tokenResp.AccessToken
	expiresIn := tokenResp.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 300
	}
	s.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)

	return s.token, nil
}

// InitiateCollection sends a collection request to ioTec Pay
func (s *IoTecService) InitiateCollection(ctx context.Context, req IoTecCollectionRequest) (*IoTecCollectionResponse, error) {
	if req.Payer == "" {
		return nil, errors.New("payer phone number is required")
	}
	if req.Amount < 500 {
		return nil, errors.New("amount must be at least UGX 500")
	}
	if req.Category == "" {
		req.Category = "MobileMoney"
	}
	if req.Currency == "" {
		req.Currency = "UGX"
	}
	if req.WalletID == "" {
		req.WalletID = s.cfg.WalletID
	}

	// Normalize Uganda phone numbers (e.g. +256770000000 -> 0770000000 or keep standard)
	cleanPhone := strings.TrimSpace(req.Payer)
	cleanPhone = strings.ReplaceAll(cleanPhone, " ", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
	if strings.HasPrefix(cleanPhone, "+256") {
		cleanPhone = "0" + cleanPhone[4:]
	} else if strings.HasPrefix(cleanPhone, "256") {
		cleanPhone = "0" + cleanPhone[3:]
	}
	req.Payer = cleanPhone

	// Check for offline simulation mode only if credentials are unset
	if s.cfg.ClientID == "" || s.cfg.ClientSecret == "" {
		return s.simulateCollection(req)
	}

	token, err := s.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	jsonPayload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("%s/api/collections/collect", s.cfg.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("iotec collection network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("iotec collection error (%d): %s", resp.StatusCode, string(body))
	}

	var res IoTecCollectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode iotec collection response: %w", err)
	}

	return &res, nil
}

// GetStatus queries the status of a transaction ID
func (s *IoTecService) GetStatus(ctx context.Context, requestID string, phone string) (*IoTecCollectionResponse, error) {
	if requestID == "" {
		return nil, errors.New("request ID is required")
	}

	// Simulation only for offline fallback or simulated ids
	if s.cfg.ClientID == "" || s.cfg.ClientSecret == "" || strings.HasPrefix(requestID, "sim-") {
		return s.simulateStatus(requestID, phone)
	}

	token, err := s.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("%s/api/collections/status/%s", s.cfg.BaseURL, requestID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("iotec status check error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("iotec status response (%d): %s", resp.StatusCode, string(body))
	}

	var res IoTecCollectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode iotec status response: %w", err)
	}

	return &res, nil
}

// Simulation logic following the official test table from api-1.json
func (s *IoTecService) simulateCollection(req IoTecCollectionRequest) (*IoTecCollectionResponse, error) {
	simID := fmt.Sprintf("sim-%d", time.Now().UnixNano())
	status := "Pending"
	msg := "USSD prompt sent to payer mobile phone"

	// Immediate failure test number
	if strings.HasPrefix(req.Payer, "011177799") {
		status = "Failed"
		msg = "Transaction failed or declined by customer"
	}

	return &IoTecCollectionResponse{
		ID:             simID,
		CreatedAt:      time.Now().UTC().Format(time.RFC3339),
		Category:       "MobileMoney",
		Status:         status,
		PaymentChannel: "MTN",
		StatusCode:     strings.ToLower(status),
		StatusMessage:  msg,
		ExternalID:     req.ExternalID,
		Amount:         req.Amount,
		Currency:       req.Currency,
	}, nil
}

func (s *IoTecService) simulateStatus(requestID string, phone string) (*IoTecCollectionResponse, error) {
	status := "Success"
	msg := "Payment confirmed successfully"

	if strings.HasPrefix(phone, "011177799") {
		status = "Failed"
		msg = "Transaction declined by customer or insufficient funds"
	} else if strings.HasPrefix(phone, "011177778") {
		status = "Pending"
		msg = "Awaiting PIN confirmation on customer device"
	}

	return &IoTecCollectionResponse{
		ID:             requestID,
		CreatedAt:      time.Now().UTC().Format(time.RFC3339),
		Category:       "MobileMoney",
		Status:         status,
		PaymentChannel: "MobileMoney",
		StatusCode:     strings.ToLower(status),
		StatusMessage:  msg,
		ExternalID:     "SIM-EXT",
		Amount:         0,
		Currency:       "UGX",
	}, nil
}
