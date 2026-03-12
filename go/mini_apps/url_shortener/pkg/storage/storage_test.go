package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestRepository_Load_FileMissing(t *testing.T) {
	repo := Repository{DbPath: "missing.json"}

	_, err := repo.Load()

	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRepository_Load_ReturnsPairs(t *testing.T) {
	urls := []URLPair{
		{"https://google.com", "abc"},
		{"https://youtube.com", "def"},
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "db.json")

	data, _ := json.Marshal(urls)
	os.WriteFile(path, data, 0644)

	repo := Repository{DbPath: path}

	got, err := repo.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(got, urls) {
		t.Fatalf("got %#v want %#v", got, urls)
	}
}

func TestRepository_Load_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "db.json")

	os.WriteFile(path, []byte("{invalid json"), 0644)

	repo := Repository{DbPath: path}

	_, err := repo.Load()

	if err == nil {
		t.Fatal("expected error for invalid json")
	}
}

func TestRepository_Store_WritesJSON(t *testing.T) {
	urls := []URLPair{
		{"https://google.com", "abc"},
		{"https://youtube.com", "def"},
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "db.json")

	repo := Repository{DbPath: path}

	err := repo.Store(urls)
	if err != nil {
		t.Fatalf("store failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read stored file: %v", err)
	}

	var got []URLPair
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("invalid json written: %v", err)
	}

	if !reflect.DeepEqual(got, urls) {
		t.Fatalf("stored data mismatch\n got: %#v\nwant: %#v", got, urls)
	}
}

func TestRepository_Store_WriteError(t *testing.T) {
	repo := Repository{DbPath: "/invalid/path/db.json"}

	err := repo.Store([]URLPair{{"a", "b"}})

	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestStorage_StoreLoad_RoundTrip(t *testing.T) {

	lstUrls := []URLPair{
		{
			LongUrl:  "https://google.com",
			ShortUrl: "https://trash-url.com/1234567890",
		},
		{
			LongUrl:  "https://youtube.com",
			ShortUrl: "https://trash-url.com/0987654321",
		},
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "db.json")

	storage := Repository{DbPath: path}

	err := storage.Store(lstUrls)

	if err != nil {
		t.Fatalf("failed to store: %v", err)
	}

	// validate stored data
	data, err := storage.Load()
	if err != nil {
		t.Fatalf("unexpected error loading from file: %v", err)
	}

	if !reflect.DeepEqual(data, lstUrls) {
		t.Fatalf("got %+v, want %+v", data, lstUrls)
	}

	if len(data) != len(lstUrls) {
		t.Fatalf("length mismatch: got %d want %d", len(data), len(lstUrls))
	}
}
