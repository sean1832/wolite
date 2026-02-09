package auth

import (
	"crypto/rand"
	"math/big"
	"os"
)

func GenerateToken(length int) ([]byte, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, length)
	for i := range token {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return []byte{}, err
		}
		token[i] = charset[randomIndex.Int64()]
	}
	return token, nil
}

func TokenExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
