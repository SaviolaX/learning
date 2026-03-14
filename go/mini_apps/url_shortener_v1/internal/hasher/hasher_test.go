package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestSha256(t *testing.T) {
	t.Run("invalid maxLen", func(t *testing.T) {
		incorrectMaxLen := 100
		_, err := Sha256("asdf", incorrectMaxLen)
		want := "invalid maxLen: 100"

		if err.Error() != want {
			t.Fatalf("expected 'invalid maxLen', got %v", err)
		}
	})

	t.Run("empty url", func(t *testing.T) {
		_, err := Sha256("", 10)

		want := "url is empty"

		if err.Error() != want {
			t.Fatalf("expected 'url is empty', got %v", err)
		}
	})

	t.Run("valid hash", func(t *testing.T) {
		url := "https://google.com"

		got, err := Sha256(url, 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		hash := sha256.Sum256([]byte(url))
		want := hex.EncodeToString(hash[:])[:10]

		if got != want {
			t.Fatalf("got %s want %s", got, want)
		}
	})
}
