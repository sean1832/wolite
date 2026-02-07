package auth

import (
	"crypto/rand"
	"fmt"

	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@)!()"
	b := make([]byte, length)
	for i := range b {
		_, err := rand.Read(b[i : i+1])
		if err != nil {
			return "", fmt.Errorf("failed to generate random string: %w", err)
		}
	}
	return string(b), nil
}

// GenerationOTPSecret generates a new OTP secret for a user and returns the secret and provisioning URL
func GenerateOTPSecret(username string) (secret string, url string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Wolite",
		AccountName: username,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to generate OTP secret: %w", err)
	}

	// Save 'secret' in the user record for later verification
	// Send 'url' to the user to generate QR code
	return key.Secret(), key.URL(), nil
}

// Validate2FA validates the provided passcode against the user's secret
func Validate2FA(passcode string, userSecret string) bool {
	// Validates the code against the secret and current time
	return totp.Validate(passcode, userSecret)
}
