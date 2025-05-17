package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"time"
)

// GenerateToken creates a unique, secure token string for email verification links.
// It combines random data with a timestamp to ensure uniqueness.
// Returns the generated token and any error encountered during creation.
func GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(b)

	timestamp := time.Now().UnixNano()
	tokenWithTimestamp := base64.URLEncoding.EncodeToString([]byte(token + strconv.FormatInt(timestamp, 10)))

	return tokenWithTimestamp, nil
}
