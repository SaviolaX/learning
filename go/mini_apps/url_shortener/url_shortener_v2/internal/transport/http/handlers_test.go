package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	application "urlShortenerV2/internal/application/url"
	dom "urlShortenerV2/internal/domain/url"
	jsonRepo "urlShortenerV2/internal/infrastructure/databases/json"
	"urlShortenerV2/internal/infrastructure/hasher"

	"github.com/stretchr/testify/assert"
)

var port = ":8080"
var host = "localhost"
var respHostUrl = "http://" + host + port + "/r/"

func TestRedirectUrl_NotFound(t *testing.T) {
	indexPath := "../../../templates/index.html"

	repo := newTestRepo(t)

	service := application.NewService(repo, hasher.Sha256, 10)

	handler := NewHandler(indexPath, service, host, port)
	s := NewServer(handler)

	req := httptest.NewRequest(http.MethodGet, "/r/qwery", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()

	s.h.RedirectUrl(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestRedirectUrl_Correct(t *testing.T) {
	indexPath := "../../../templates/index.html"

	repo := newTestRepo(t)

	pair := dom.UrlPair{
		Code:      "qwery",
		OriginUrl: dom.Url{Value: "https://google.com"},
	}

	err := repo.Save(context.Background(), pair)

	assert.NoError(t, err)

	service := application.NewService(repo, hasher.Sha256, 10)

	handler := NewHandler(indexPath, service, host, port)
	s := NewServer(handler)

	req := httptest.NewRequest(http.MethodGet, "/r/qwery", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()

	s.h.RedirectUrl(rec, req)

	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, "https://google.com", rec.Header().Get("Location"))
}

func TestRedirectUrl_MethodNotAllowed(t *testing.T) {
	indexPath := "../../../templates/index.html"

	repo := newTestRepo(t)
	service := application.NewService(repo, hasher.Sha256, 10)

	handler := NewHandler(indexPath, service, host, port)
	s := NewServer(handler)

	req := httptest.NewRequest(http.MethodPost, "/r/1234567890", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()

	s.h.RedirectUrl(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestShortenUrl_BadRequest(t *testing.T) {
	indexPath := "../../../templates/index.html"

	repo := newTestRepo(t)
	service := application.NewService(repo, hasher.Sha256, 10)

	handler := NewHandler(indexPath, service, host, port)
	s := NewServer(handler)

	form := url.Values{}
	form.Add("url", "")

	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()

	s.h.ShortenUrl(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestShortenUrl_Correct(t *testing.T) {
	indexPath := "../../../templates/index.html"

	repo := newTestRepo(t)
	service := application.NewService(repo, hasher.Sha256, 10)

	handler := NewHandler(indexPath, service, host, port)
	s := NewServer(handler)

	form := url.Values{}
	form.Add("url", "google.com")

	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()

	s.h.ShortenUrl(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), respHostUrl)
}

func TestShortenUrl_MethodNotAllowed(t *testing.T) {
	indexPath := "../../../templates/index.html"

	repo := newTestRepo(t)
	service := application.NewService(repo, hasher.Sha256, 10)

	handler := NewHandler(indexPath, service, host, port)
	s := NewServer(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()

	s.h.ShortenUrl(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestHomePage_GET(t *testing.T) {
	indexPath := "../../../templates/index.html"

	repo := newTestRepo(t)
	service := application.NewService(repo, hasher.Sha256, 10)

	handler := NewHandler(indexPath, service, host, port)
	s := NewServer(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()

	s.h.homePage(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "<html")
}

func TestHomePage_MethodNotAllowed(t *testing.T) {
	indexPath := "../../../templates/index.html"

	repo := newTestRepo(t)
	service := application.NewService(repo, hasher.Sha256, 10)

	handler := NewHandler(indexPath, service, host, port)
	s := NewServer(handler)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()

	s.h.homePage(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func newTestRepo(t *testing.T) *jsonRepo.Repository {
	t.Helper()
	dir := t.TempDir()
	return jsonRepo.NewRepository(filepath.Join(dir, "db.json"))
}
