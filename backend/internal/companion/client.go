package companion

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client is a client for the Companion API.
type Client struct {
	BaseURL     string
	Token       string
	Fingerprint string // Expected SHA-256 fingerprint (hex encoded)
	HTTPClient  *http.Client
}

// NewClient creates a new Companion API client.
// It configures the HTTP client to strictly verify the peer certificate against the provided fingerprint.
func NewClient(baseURL, token, fingerprint string) (*Client, error) {
	if baseURL == "" || token == "" {
		return nil, errors.New("base URL and token are required")
	}

	// Normalize fingerprint
	fingerprint = strings.ToLower(strings.ReplaceAll(fingerprint, ":", ""))

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // We verify manually in VerifyConnection
			VerifyConnection: func(cs tls.ConnectionState) error {
				if len(cs.PeerCertificates) == 0 {
					return errors.New("no certificate presented by peer")
				}
				leaf := cs.PeerCertificates[0]

				// Calculate fingerprint
				hash := sha256.Sum256(leaf.Raw)
				actualFingerprint := hex.EncodeToString(hash[:])

				if actualFingerprint != fingerprint {
					return fmt.Errorf("certificate fingerprint mismatch: expected %s, got %s", fingerprint, actualFingerprint)
				}
				return nil
			},
		},
	}

	return &Client{
		BaseURL:     strings.TrimRight(baseURL, "/"),
		Token:       token,
		Fingerprint: fingerprint,
		HTTPClient: &http.Client{
			Timeout:   5 * time.Second,
			Transport: transport,
		},
	}, nil
}

// GetFingerprint connects to the given URL (insecurely) and retrieves the leaf certificate's SHA-256 fingerprint.
// This implements the Trust-On-First-Use (TOFU) flow.
func GetFingerprint(ctx context.Context, urlStr string) (string, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Intentionally insecure to fetch the cert
		},
	}
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	// We use the health endpoint for a lightweight check
	req, err := http.NewRequestWithContext(ctx, "GET", strings.TrimRight(urlStr, "/")+"/api/v1/health", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.TLS == nil || len(resp.TLS.PeerCertificates) == 0 {
		return "", errors.New("no TLS state or certificates found")
	}

	leaf := resp.TLS.PeerCertificates[0]
	hash := sha256.Sum256(leaf.Raw)
	return hex.EncodeToString(hash[:]), nil
}

// Ping checks if the companion is reachable and healthy.
func (c *Client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/api/v1/health", nil)
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("companion returned status: %d", resp.StatusCode)
	}
	return nil
}

// PowerAction represents a power command.
type PowerAction string

const (
	ActionShutdown  PowerAction = "shutdown"
	ActionReboot    PowerAction = "reboot"
	ActionSleep     PowerAction = "sleep"
	ActionHibernate PowerAction = "hibernate"
)

// Power sends a power command to the companion.
func (c *Client) Power(ctx context.Context, action PowerAction) error {
	url := fmt.Sprintf("%s/api/v1/%s", c.BaseURL, action)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("companion error (%d): %s", resp.StatusCode, string(body))
	}
	return nil
}
