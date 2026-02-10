package tls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

// GenerateSelfSignedCert creates a new self-signed certificate and private key
func GenerateSelfSignedCert(certPath, keyPath string) error {
	// generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Create certificate template

	// Generate a 128-bit serial number.
	// Use bit shifting (1 << 128) to define the upper bound [0, 2^128).
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, max)
	if err != nil {
		return err
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"sean1832/wolite"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(2, 0, 0), // expire after 2 years
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},

		// support both localhost and LAN IPs
		IPAddresses: []net.IP{
			net.ParseIP("127.0.0.1"),
			net.ParseIP("::1"),
		},
		DNSNames: []string{"localhost"},
	}

	// add all local ip to certificates
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil || ipnet.IP.To16() != nil {
					template.IPAddresses = append(template.IPAddresses, ipnet.IP)
				}
			}
		}
	}

	// create self-signed certificate
	certBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template, // self-signed: issuer = subject
		&privateKey.PublicKey,
		privateKey,
	)

	if err != nil {
		return err
	}

	// write certificate to file
	certFile, err := os.Create(certPath)
	if err != nil {
		return err
	}
	defer certFile.Close()

	if err := pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return err
	}

	// write private key to file
	keyFile, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	if err := pem.Encode(keyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return err
	}

	return nil
}

// CertExists checks if certificate file already exist
func CertExists(certPath, keyPath string) bool {
	_, certErr := os.Stat(certPath)
	_, keyErr := os.Stat(keyPath)
	return certErr == nil && keyErr == nil
}

// GetCertFingerprint retrieves the SHA-256 fingerprint of the certificate
func GetCertFingerprint(certPath string) (string, error) {
	certData, err := os.ReadFile(certPath)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(certData)
	if block == nil || block.Type != "CERTIFICATE" {
		return "", fmt.Errorf("failed to decode PEM block containing certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", err
	}
	// Calculate SHA-256 hash of the Raw ASN.1 DER data
	hash := sha256.Sum256(cert.Raw)
	return fmt.Sprintf("%X", hash), nil
}
