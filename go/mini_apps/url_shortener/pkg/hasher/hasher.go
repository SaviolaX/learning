package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

type Hasher struct {
	Url string
}

func (h Hasher) Sha256(maxLen int) (string, error) {
	if len(h.Url) == 0 {
		return "", errors.New("no url passed into HashUrl")
	}

	hash := sha256.Sum256([]byte(h.Url))
	hashStr := hex.EncodeToString(hash[:])

	return hashStr[:maxLen], nil
}
