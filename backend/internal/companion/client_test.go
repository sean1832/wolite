package companion

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// generateCert generates a self-signed certificate and private key for testing.
func generateCert() (tls.Certificate, string, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, "", err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, "", err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return tls.Certificate{}, "", err
	}

	// Calculate fingerprint
	hash := sha256.Sum256(derBytes)
	fingerprint := hex.EncodeToString(hash[:])

	return cert, fingerprint, nil
}

func TestCompanionClient(t *testing.T) {
	// 1. Setup Mock Server with Self-Signed Cert
	cert, expectedFingerprint, err := generateCert()
	if err != nil {
		t.Fatalf("Failed to generate cert: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/api/v1/shutdown", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("shutting down"))
	})

	server := httptest.NewUnstartedServer(mux)
	server.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}
	server.StartTLS()
	defer server.Close()

	ctx := context.Background()

	// 2. Test GetFingerprint
	t.Run("GetFingerprint", func(t *testing.T) {
		fp, err := GetFingerprint(ctx, server.URL)
		if err != nil {
			t.Fatalf("GetFingerprint failed: %v", err)
		}
		if fp != expectedFingerprint {
			t.Errorf("Fingerprint mismatch. Expected %s, got %s", expectedFingerprint, fp)
		}
	})

	// 3. Test NewClient with Correct Fingerprint
	t.Run("Client_CorrectFingerprint", func(t *testing.T) {
		client, err := NewClient(server.URL, "test-token", expectedFingerprint)
		if err != nil {
			t.Fatalf("NewClient failed: %v", err)
		}

		// Test Ping
		if err := client.Ping(ctx); err != nil {
			t.Errorf("Ping failed: %v", err)
		}

		// Test Power
		if err := client.Power(ctx, ActionShutdown); err != nil {
			t.Errorf("Power failed: %v", err)
		}
	})

	// 4. Test NewClient with Incorrect Fingerprint
	t.Run("Client_IncorrectFingerprint", func(t *testing.T) {
		wrongFingerprint := strings.Repeat("a", 64)
		client, err := NewClient(server.URL, "test-token", wrongFingerprint)
		if err != nil {
			t.Fatalf("NewClient failed: %v", err)
		}

		// Ping should fail
		if err := client.Ping(ctx); err == nil {
			t.Error("Ping succeeded with wrong fingerprint, expected failure")
		} else if !strings.Contains(err.Error(), "fingerprint mismatch") {
			t.Errorf("Unexpected error message: %v", err)
		}
	})

	// 5. Test Invalid Token
	t.Run("Client_InvalidToken", func(t *testing.T) {
		client, err := NewClient(server.URL, "wrong-token", expectedFingerprint)
		if err != nil {
			t.Fatalf("NewClient failed: %v", err)
		}

		// Power should fail (Unauthorized)
		if err := client.Power(ctx, ActionShutdown); err == nil {
			t.Error("Power succeeded with wrong token, expected failure")
		} else if !strings.Contains(err.Error(), "401") {
			t.Errorf("Unexpected error message: %v", err)
		}
	})
}
