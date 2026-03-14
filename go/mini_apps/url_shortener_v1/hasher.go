package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

func Sha256(url string, maxLen int) (string, error) {
	if len(url) == 0 {
		return "", errors.New("url is empty")
	}

	hash := sha256.Sum256([]byte(url))
	hashStr := hex.EncodeToString(hash[:])

	if maxLen <= 0 || maxLen > len(hashStr) {
		return "", fmt.Errorf("invalid maxLen: %d", maxLen)
	}

	return hashStr[:maxLen], nil
}
