package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRefreshToken() (string, error) {
	randomByte := make([]byte, 64)

	_, err := rand.Read(randomByte)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	token := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomByte)

	return token, nil
}
