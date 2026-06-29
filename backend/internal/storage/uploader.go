package storage

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// googleDriveScope must be the full Drive scope (not drive.file): the admin
// supplies the ID of a folder that already exists in their own Drive rather
// than one created or opened by this app, and drive.file only grants access
// to files/folders the app itself created or that were explicitly opened via
// the Drive picker. Without the broader scope, uploads into a pre-existing
// folder fail with a 403 even though the OAuth credentials are valid.
const googleDriveScope = "https://www.googleapis.com/auth/drive"

// googleDriveHTTPTimeout bounds how long a single Drive API call (token
// refresh, file upload, or permission change) may take before failing.
const googleDriveHTTPTimeout = 60 * time.Second

type Scope string

const (
	ScopePublic      Scope = "public"
	ScopeApplication Scope = "application"
)

type UploadInput struct {
	Scope       Scope
	Reader      io.Reader
	Filename    string
	ContentType string
	Size        int64
	Subdir      string
}

type UploadResult struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Provider string `json:"provider"`
	Key      string `json:"key,omitempty"`
}

type Uploader struct {
	settings Settings
	client   *http.Client
}

func NewUploader(settings Settings) *Uploader {
	settings.Normalize()
	return &Uploader{settings: settings, client: &http.Client{Timeout: googleDriveHTTPTimeout}}
}

func (u *Uploader) Upload(ctx context.Context, in UploadInput) (*UploadResult, error) {
	if in.Reader == nil {
		return nil, errors.New("reader is required")
	}
	provider := u.providerFor(in.Scope)
	switch provider {
	case "local":
		return u.uploadLocal(in)
	case "google_drive":
		return u.uploadGoogleDrive(ctx, in)
	case "s3":
		return u.uploadS3(ctx, in)
	default:
		return nil, fmt.Errorf("unsupported storage provider %q", provider)
	}
}

func (u *Uploader) providerFor(scope Scope) string {
	if scope == ScopeApplication {
		return u.settings.ApplicationProvider
	}
	return u.settings.PublicProvider
}

func (u *Uploader) uploadLocal(in UploadInput) (*UploadResult, error) {
	basePath := u.settings.LocalPublicPath
	urlPrefix := u.settings.LocalPublicURLPrefix
	if in.Scope == ScopeApplication {
		basePath = u.settings.LocalAppPath
		urlPrefix = u.settings.LocalAppURLPrefix
	}
	ext := strings.ToLower(filepath.Ext(in.Filename))
	filename := uniqueName(ext)
	subdir := cleanSegment(in.Subdir)
	dir := filepath.Join(basePath, subdir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	target := filepath.Join(dir, filename)
	out, err := os.Create(target)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	if _, err := io.Copy(out, in.Reader); err != nil {
		return nil, err
	}
	return &UploadResult{
		URL:      joinURLPath(urlPrefix, subdir, filename),
		Filename: filename,
		Provider: "local",
		Key:      path.Join(subdir, filename),
	}, nil
}

func (u *Uploader) uploadGoogleDrive(ctx context.Context, in UploadInput) (*UploadResult, error) {
	ctx, cancel := context.WithTimeout(ctx, googleDriveHTTPTimeout)
	defer cancel()
	client, err := u.googleDriveClient(ctx)
	if err != nil {
		return nil, err
	}
	content, err := io.ReadAll(in.Reader)
	if err != nil {
		return nil, err
	}
	ext := strings.ToLower(filepath.Ext(in.Filename))
	filename := uniqueName(ext)
	meta := map[string]interface{}{"name": filename}
	if u.settings.GoogleDriveFolderID != "" {
		meta["parents"] = []string{u.settings.GoogleDriveFolderID}
	}
	metaBytes, _ := json.Marshal(meta)
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	metaPart, err := writer.CreatePart(map[string][]string{"Content-Type": {"application/json; charset=UTF-8"}})
	if err != nil {
		return nil, err
	}
	if _, err := metaPart.Write(metaBytes); err != nil {
		return nil, err
	}
	mediaHeader := map[string][]string{"Content-Type": {contentType(in.ContentType, filename)}}
	mediaPart, err := writer.CreatePart(mediaHeader)
	if err != nil {
		return nil, err
	}
	if _, err := mediaPart.Write(content); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://www.googleapis.com/upload/drive/v3/files?uploadType=multipart&supportsAllDrives=true&fields=id,name,webViewLink", &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("google drive upload failed: %s %s", resp.Status, strings.TrimSpace(string(b)))
	}
	var created struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		WebViewLink string `json:"webViewLink"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, err
	}
	if u.settings.GoogleDriveMakePublic {
		if err := makeDriveFilePublic(ctx, client, created.ID); err != nil {
			slog.Warn("google drive: could not make uploaded file public", "file_id", created.ID, "error", err)
		}
	}
	return &UploadResult{
		URL:      "https://drive.google.com/uc?id=" + url.QueryEscape(created.ID),
		Filename: filename,
		Provider: "google_drive",
		Key:      created.ID,
	}, nil
}

// googleDriveOAuthConfig builds the OAuth2 config shared by the connect,
// exchange, and authenticated-client code paths so the client ID/secret,
// endpoint, and scope can never drift between them.
func googleDriveOAuthConfig(clientID, clientSecret, redirectURI string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  redirectURI,
		Scopes:       []string{googleDriveScope},
	}
}

// GoogleDriveAuthURL builds the Google consent-screen URL for the admin to
// authorize this app against their Drive account. access_type=offline plus
// prompt=consent guarantees a refresh token comes back even on a reconnect.
func GoogleDriveAuthURL(clientID, redirectURI, state string) string {
	cfg := googleDriveOAuthConfig(clientID, "", redirectURI)
	return cfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
}

// ExchangeGoogleDriveCode trades the one-time authorization code from the
// OAuth callback for a refresh token, which is what authenticates every
// subsequent upload — the admin never has to see or paste this value.
func ExchangeGoogleDriveCode(ctx context.Context, clientID, clientSecret, redirectURI, code string) (string, error) {
	cfg := googleDriveOAuthConfig(clientID, clientSecret, redirectURI)
	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("google drive token exchange failed: %w", err)
	}
	if strings.TrimSpace(token.RefreshToken) == "" {
		return "", errors.New("google did not return a refresh token; remove this app's access in your Google account's third-party access settings and connect again")
	}
	return token.RefreshToken, nil
}

func (u *Uploader) googleDriveClient(ctx context.Context) (*http.Client, error) {
	if u.settings.GoogleDriveClientID == "" || u.settings.GoogleDriveClientSecret == "" {
		return nil, errors.New("google drive client id and client secret are not configured")
	}
	if u.settings.GoogleDriveRefreshToken == "" {
		return nil, errors.New("google drive is not connected: save the client id, client secret and folder id, then connect your Google account")
	}
	cfg := googleDriveOAuthConfig(u.settings.GoogleDriveClientID, u.settings.GoogleDriveClientSecret, "")
	token := &oauth2.Token{RefreshToken: u.settings.GoogleDriveRefreshToken}
	return cfg.Client(ctx, token), nil
}

func makeDriveFilePublic(ctx context.Context, client *http.Client, fileID string) error {
	body := strings.NewReader(`{"role":"reader","type":"anyone"}`)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://www.googleapis.com/drive/v3/files/"+url.PathEscape(fileID)+"/permissions?supportsAllDrives=true", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("google drive permission failed: %s", resp.Status)
	}
	return nil
}

func (u *Uploader) uploadS3(ctx context.Context, in UploadInput) (*UploadResult, error) {
	if u.settings.S3Bucket == "" || u.settings.S3AccessKeyID == "" || u.settings.S3SecretAccessKey == "" {
		return nil, errors.New("s3 bucket and credentials are not configured")
	}
	content, err := io.ReadAll(in.Reader)
	if err != nil {
		return nil, err
	}
	ext := strings.ToLower(filepath.Ext(in.Filename))
	filename := uniqueName(ext)
	key := path.Join(cleanSegment(u.settings.S3Prefix), cleanSegment(in.Subdir), filename)
	endpoint, err := u.s3ObjectURL(key)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	payloadHash := sha256Hex(content)
	req.Header.Set("Content-Type", contentType(in.ContentType, filename))
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)
	req.Header.Set("X-Amz-Date", now.Format("20060102T150405Z"))
	signS3(req, u.settings, now, payloadHash)
	resp, err := u.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("s3 upload failed: %s %s", resp.Status, strings.TrimSpace(string(b)))
	}
	return &UploadResult{URL: u.s3PublicURL(key), Filename: filename, Provider: "s3", Key: key}, nil
}

func (u *Uploader) s3ObjectURL(key string) (string, error) {
	if u.settings.S3Endpoint != "" {
		base := strings.TrimRight(u.settings.S3Endpoint, "/")
		if u.settings.S3ForcePathStyle {
			return base + "/" + url.PathEscape(u.settings.S3Bucket) + "/" + escapeKey(key), nil
		}
		parsed, err := url.Parse(base)
		if err != nil {
			return "", err
		}
		parsed.Host = u.settings.S3Bucket + "." + parsed.Host
		parsed.Path = "/" + escapeKey(key)
		return parsed.String(), nil
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.settings.S3Bucket, u.settings.S3Region, escapeKey(key)), nil
}

func (u *Uploader) s3PublicURL(key string) string {
	if u.settings.S3PublicBaseURL != "" {
		return strings.TrimRight(u.settings.S3PublicBaseURL, "/") + "/" + escapeKey(key)
	}
	endpoint, _ := u.s3ObjectURL(key)
	return endpoint
}

func signS3(req *http.Request, cfg Settings, now time.Time, payloadHash string) {
	date := now.Format("20060102")
	credentialScope := date + "/" + cfg.S3Region + "/s3/aws4_request"
	headers := signedHeaders(req)
	canonical := req.Method + "\n" +
		req.URL.EscapedPath() + "\n" +
		canonicalQuery(req.URL.Query()) + "\n" +
		canonicalHeaders(req, headers) + "\n" +
		strings.Join(headers, ";") + "\n" +
		payloadHash
	stringToSign := "AWS4-HMAC-SHA256\n" +
		now.Format("20060102T150405Z") + "\n" +
		credentialScope + "\n" +
		sha256Hex([]byte(canonical))
	key := s3SigningKey(cfg.S3SecretAccessKey, date, cfg.S3Region)
	signature := hex.EncodeToString(hmacSHA256(key, stringToSign))
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential="+cfg.S3AccessKeyID+"/"+credentialScope+", SignedHeaders="+strings.Join(headers, ";")+", Signature="+signature)
}

func signedHeaders(req *http.Request) []string {
	headers := []string{"host"}
	for k := range req.Header {
		headers = append(headers, strings.ToLower(k))
	}
	sort.Strings(headers)
	return compact(headers)
}

func canonicalHeaders(req *http.Request, headers []string) string {
	var b strings.Builder
	for _, h := range headers {
		value := req.URL.Host
		if h != "host" {
			value = strings.Join(req.Header.Values(http.CanonicalHeaderKey(h)), ",")
		}
		b.WriteString(h)
		b.WriteByte(':')
		b.WriteString(strings.Join(strings.Fields(value), " "))
		b.WriteByte('\n')
	}
	return b.String()
}

func canonicalQuery(values url.Values) string {
	if len(values) == 0 {
		return ""
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0)
	for _, key := range keys {
		vals := append([]string(nil), values[key]...)
		sort.Strings(vals)
		for _, value := range vals {
			parts = append(parts, url.QueryEscape(key)+"="+url.QueryEscape(value))
		}
	}
	return strings.Join(parts, "&")
}

func s3SigningKey(secret, date, region string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secret), date)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, "s3")
	return hmacSHA256(kService, "aws4_request")
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	_, _ = h.Write([]byte(data))
	return h.Sum(nil)
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func uniqueName(ext string) string {
	return fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), randomHex(8), ext)
}

func randomHex(n int) string {
	b := make([]byte, n)
	sum := sha256.Sum256([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	copy(b, sum[:n])
	return hex.EncodeToString(b)
}

func cleanSegment(value string) string {
	value = strings.Trim(strings.ReplaceAll(value, "\\", "/"), "/")
	if value == "" || value == "." {
		return ""
	}
	parts := strings.Split(value, "/")
	cleaned := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "." || part == ".." {
			continue
		}
		cleaned = append(cleaned, part)
	}
	return path.Join(cleaned...)
}

func joinURLPath(prefix string, parts ...string) string {
	all := []string{strings.TrimRight(prefix, "/")}
	for _, part := range parts {
		if part = strings.Trim(part, "/"); part != "" {
			all = append(all, part)
		}
	}
	return strings.Join(all, "/")
}

func contentType(given, filename string) string {
	if given != "" {
		return given
	}
	if t := mime.TypeByExtension(filepath.Ext(filename)); t != "" {
		return t
	}
	return "application/octet-stream"
}

func escapeKey(key string) string {
	parts := strings.Split(key, "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}

func compact(values []string) []string {
	out := values[:0]
	var prev string
	for _, value := range values {
		if value == prev {
			continue
		}
		out = append(out, value)
		prev = value
	}
	return out
}
