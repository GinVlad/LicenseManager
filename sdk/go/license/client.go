// Package license provides a client for validating licenses against LicenseManager.
// Import this in any Wails (or CLI) app to validate licenses with a single function call.
package license

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// Config holds the settings for the license client.
type Config struct {
	ServerURL  string // e.g. "https://licenses.yourdomain.com"
	AppSlug    string // e.g. "ebay-creator"
	AppAPIKey  string // X-App-Key header value (from admin panel)
	DeviceName string // optional: e.g. "Home PC"
}

// LicenseInfo is the result of a successful license validation.
type LicenseInfo struct {
	Valid      bool
	Key        string
	ExpiresAt  time.Time
	MaxThreads int
	HWID       string
}

// Client validates licenses against the LicenseManager server.
type Client struct {
	cfg    Config
	hwid   string
	http   *http.Client
}

// New creates a new license client.
func New(cfg Config) *Client {
	return &Client{
		cfg:  cfg,
		hwid: generateHWID(),
		http: &http.Client{Timeout: 15 * time.Second},
	}
}

// HWID returns this machine's hardware ID.
func (c *Client) HWID() string {
	return c.hwid
}

// Validate sends the license key and HWID to the server.
// Returns LicenseInfo on success, or an error describing the failure.
func (c *Client) Validate(key string) (*LicenseInfo, error) {
	payload := map[string]string{
		"appSlug":    c.cfg.AppSlug,
		"key":        key,
		"hwid":       c.hwid,
		"deviceName": c.cfg.DeviceName,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", c.cfg.ServerURL+"/api/v1/license/validate", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Key", c.cfg.AppAPIKey)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			Valid      bool       `json:"valid"`
			Error      string     `json:"error"`
			ExpiresAt  *time.Time `json:"expiresAt"`
			MaxThreads int        `json:"maxThreads"`
		} `json:"data"`
		Error string `json:"error"`
	}

	data, _ := io.ReadAll(res.Body)
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if !resp.Success {
		return nil, fmt.Errorf("server error: %s", resp.Error)
	}
	if !resp.Data.Valid {
		return nil, fmt.Errorf("%s", resp.Data.Error)
	}

	info := &LicenseInfo{
		Valid:      true,
		Key:        key,
		MaxThreads: resp.Data.MaxThreads,
		HWID:       c.hwid,
	}
	if resp.Data.ExpiresAt != nil {
		info.ExpiresAt = *resp.Data.ExpiresAt
	}
	return info, nil
}

// generateHWID produces a stable 32-char machine fingerprint.
func generateHWID() string {
	var parts []string

	ifaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range ifaces {
			if len(iface.HardwareAddr) > 0 {
				parts = append(parts, iface.HardwareAddr.String())
				break
			}
		}
	}

	hostname, _ := os.Hostname()
	parts = append(parts, hostname, os.Getenv("USER"))

	hash := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(hash[:16])
}
