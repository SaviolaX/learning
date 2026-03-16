package json

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestRepository_GetByCode_LoadError(t *testing.T) {
	hashCode := "asdf"

	repo := Repository{DbPath: ""}

	_, got := repo.GetByCode(hashCode)
	want := "url not found"
	if got.Error() != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRepository_GetByCode_NotFound(t *testing.T) {
	hashCode := "asdf"

	dir := t.TempDir()
	path := filepath.Join(dir, "db.json")

	repo := Repository{DbPath: path}

	_, got := repo.GetByCode(hashCode)
	want := "url not found"

	if got.Error() != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRepository_GetByCode_Found(t *testing.T) {
	hashCode := "asdf"
	url := "google.com"
	urls := UrlPair{Code: hashCode, OriginalUrl: url}

	dir := t.TempDir()
	path := filepath.Join(dir, "db.json")

	repo := Repository{DbPath: path}

	if err := repo.Store(urls); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := repo.GetByCode(hashCode)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := UrlPair{Code: hashCode, OriginalUrl: url}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}

}

func TestRepository_Load_FileMissing(t *testing.T) {
	repo := Repository{DbPath: "missing.json"}

	got, err := repo.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var want []UrlPair

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}

}

func TestRepository_Load_ReturnsPairs(t *testing.T) {
	urls := []UrlPair{
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
	urls := []UrlPair{
		{"asdf", "https://youtube.com"},
	}

	url := UrlPair{Code: "asdf", OriginalUrl: "https://youtube.com"}

	dir := t.TempDir()
	path := filepath.Join(dir, "db.json")

	repo := Repository{DbPath: path}

	err := repo.Store(url)
	if err != nil {
		t.Fatalf("store failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read stored file: %v", err)
	}

	var got []UrlPair
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("invalid json written: %v", err)
	}

	if !reflect.DeepEqual(got, urls) {
		t.Fatalf("stored data mismatch\n got: %#v\nwant: %#v", got, urls)
	}
}

func TestRepository_Store_WriteError(t *testing.T) {
	repo := Repository{DbPath: "/invalid/path/db.json"}

	err := repo.Store(UrlPair{Code: "a", OriginalUrl: "b"})

	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestStorage_StoreLoad_RoundTrip(t *testing.T) {

	lstUrls := []UrlPair{
		{
			OriginalUrl: "https://youtube.com",
			Code:        "asdf",
		},
	}

	url := UrlPair{Code: "asdf", OriginalUrl: "https://youtube.com"}

	dir := t.TempDir()
	path := filepath.Join(dir, "db.json")

	storage := Repository{DbPath: path}

	err := storage.Store(url)

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
