package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomString генерирует случайную hex-строку указанной длины
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
