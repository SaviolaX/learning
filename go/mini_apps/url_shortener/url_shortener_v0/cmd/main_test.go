package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"urlShortener/pkg/storage"
)

func TestShortenUrl_MethodNotAllowed(t *testing.T) {
	repo := storage.Repository{DbPath: ""}

	s := &Server{
		repo:      repo,
		indexPath: "../templates/index.html",
	}

	router := http.NewServeMux()
	router.HandleFunc("/short-url", s.shortenUrl)

	req := httptest.NewRequest(http.MethodGet, "/short-url", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d got %d", http.StatusMethodNotAllowed, rec.Code)
	}

}

func TestHomePage_MethodNotAllowed(t *testing.T) {
	repo := storage.Repository{DbPath: ""}

	s := &Server{
		repo:      repo,
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
	repo := storage.Repository{DbPath: ""}

	s := &Server{
		repo:      repo,
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

func TestBuildMap_MapLoadedData(t *testing.T) {
	pairs := []storage.URLPair{
		{
			LongUrl:  "asdf",
			ShortUrl: "qwe",
		},

		{
			LongUrl:  "asdf",
			ShortUrl: "qwe",
		},
	}

	got := BuildMap(pairs)
	want := map[string]string{
		pairs[0].ShortUrl: pairs[0].LongUrl,
		pairs[1].ShortUrl: pairs[1].LongUrl,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestBuildMap_MapEmptyList(t *testing.T) {
	pairs := []storage.URLPair{}

	got := BuildMap(pairs)
	want := map[string]string{}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}
