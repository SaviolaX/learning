package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"testing/iotest"
	"urlShortenerV1/internal/hasher"
	jsonRepo "urlShortenerV1/internal/repositorie/json"

	"github.com/stretchr/testify/assert"
)

const respHostUrl = "http://localhost:8080/r/"
const testUrl = "https://google.com"

func TestRedirectUrl_Redirects(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "db.json")
	repo := jsonRepo.Repository{DbPath: path}

	err := repo.Store(jsonRepo.UrlPair{
		Code:        "abc123",
		OriginalUrl: "https://google.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := Server{db: &repo}

	req := httptest.NewRequest(http.MethodGet, "/r/abc123", nil)
	rec := httptest.NewRecorder()

	s.RedirectUrl(rec, req)

	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "https://google.com", rec.Header().Get("Location"))
}

func TestRedirectUrl_UrlNotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "db.json")

	repo := jsonRepo.Repository{DbPath: path}

	s := Server{
		db:        &repo,
		indexPath: "../templates/index.html",
	}

	req := httptest.NewRequest(http.MethodGet, "/r/123456679", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	s.RedirectUrl(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestRedirectUrl_MethodNotAllowed(t *testing.T) {
	repo := jsonRepo.Repository{DbPath: ""}

	s := Server{
		db:        &repo,
		indexPath: "../templates/index.html",
	}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	s.RedirectUrl(rec, req)

	assert.Equal(t, rec.Code, http.StatusMethodNotAllowed)
}

func TestShortenUrl_GetByCodeError(t *testing.T) {

	repo := jsonRepo.Repository{DbPath: ""}
	s := Server{
		db:        &repo,
		indexPath: "../templates/index.html",
	}

	form := url.Values{}
	form.Add("url", "google.com")

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	s.ShortenUrl(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestShortenUrl_ReturnExistingUrl(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "db.json")

	repo := jsonRepo.Repository{DbPath: path}
	hashCode, err := hasher.Sha256(testUrl, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := repo.Store(jsonRepo.UrlPair{Code: hashCode, OriginalUrl: testUrl}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := Server{
		db:        &repo,
		indexPath: "../templates/index.html",
	}

	form := url.Values{}
	form.Add("url", "google.com")

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	s.ShortenUrl(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), respHostUrl)

	pairs, err := repo.Load()
	assert.NoError(t, err)
	assert.Len(t, pairs, 1)
	assert.Equal(t, "https://google.com", pairs[0].OriginalUrl)

}

func TestShortenUrl_ValidationError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "db.json")

	repo := jsonRepo.Repository{DbPath: path}
	s := Server{
		db:        &repo,
		indexPath: "../templates/index.html",
	}

	form := url.Values{}
	form.Add("url", "!asdf")

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))

	rec := httptest.NewRecorder()
	s.ShortenUrl(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestShortenUrl_CannotParseForm(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "db.json")

	repo := jsonRepo.Repository{DbPath: path}
	s := Server{
		db:        &repo,
		indexPath: "../templates/index.html",
	}

	req := httptest.NewRequest(http.MethodPost, "/", iotest.ErrReader(errors.New("read error")))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	s.ShortenUrl(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestShortenUrl_ShortUrlCreated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "db.json")

	repo := jsonRepo.Repository{DbPath: path}
	s := Server{
		db:        &repo,
		indexPath: "../templates/index.html",
	}

	form := url.Values{}
	form.Add("url", "google.com")

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	s.ShortenUrl(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), respHostUrl)

	pairs, err := repo.Load()
	assert.NoError(t, err)
	assert.Len(t, pairs, 1)
	assert.Equal(t, "https://google.com", pairs[0].OriginalUrl)

}

func TestShortenUrl_MethodNotAllowed(t *testing.T) {
	repo := jsonRepo.Repository{DbPath: ""}

	s := Server{
		db:        &repo,
		indexPath: "../templates/index.html",
	}

	router := http.NewServeMux()
	router.HandleFunc("/", s.ShortenUrl)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, rec.Code, http.StatusMethodNotAllowed)
}

func TestHomePage_MethodNotAllowed(t *testing.T) {
	repo := jsonRepo.Repository{DbPath: ""}

	s := Server{
		db:        &repo,
		indexPath: "../templates/index.html",
	}

	router := http.NewServeMux()
	router.HandleFunc("/", s.homePage)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestHomePage_GET(t *testing.T) {
	repo := jsonRepo.Repository{DbPath: ""}

	s := Server{
		db:        &repo,
		indexPath: "../templates/index.html",
	}

	router := http.NewServeMux()
	router.HandleFunc("/", s.homePage)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d got %d", http.StatusOK, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "<html") {
		t.Fatalf("expected html response")
	}
}
