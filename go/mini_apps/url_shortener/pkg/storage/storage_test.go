package storage

import (
	"encoding/json"
	"os"
	"testing"
)

func TestStorage_Store_WriteValidJSON(t *testing.T) {

	longUrl := "https://google.com"
	shortUrl := "https://trash-url.com"

	urlsPair := URLPair{
		LongUrl:  longUrl,
		ShortUrl: shortUrl,
	}

	// create tmp json file
	tmpFile, err := os.CreateTemp("", "store_test_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// store into the file
	err = urlsPair.Store(tmpFile.Name())
	if err != nil {
		t.Errorf("store failed: %v", err)
	}

	// read data from file
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("cannot read file: %v", err)
	}

	var result URLPair
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if result != urlsPair {
		t.Fatalf("stored data mismatch")
	}

}
