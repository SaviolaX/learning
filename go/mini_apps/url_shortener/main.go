package main

import (
	"fmt"
	"log"
	"net/http"
	hasher "urlShortener/pkg/hasher"
	storage "urlShortener/pkg/storage"
)

const (
	hashLen        = 10
	defaultUrl     = "http://localhost:8080/r/"
	defaultStorage = "urlShortenerDB.json"
)

type Server struct {
	repo storage.Repository
}

func main() {

	repo := storage.Repository{DbPath: defaultStorage}

	s := &Server{repo: repo}

	router := http.NewServeMux()

	router.HandleFunc("/", s.homePage)
	router.HandleFunc("/short-url", s.shortenUrl)
	router.HandleFunc("/r/", s.redirectUrl)

	port := ":8080"
	fmt.Println("Starting server on port", port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func buildMap(pairs []storage.URLPair) map[string]string {
	m := make(map[string]string)

	for _, p := range pairs {
		m[p.ShortUrl] = p.LongUrl
	}

	return m
}

func (s *Server) homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./templates/index.html")
}

func (s *Server) redirectUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	shortUrl := "http://" + r.Host + r.URL.Path

	data, err := s.repo.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	urlMap := buildMap(data)
	longUrl, ok := urlMap[shortUrl]
	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longUrl, http.StatusFound)
}

func (s *Server) shortenUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "cannot parse form", http.StatusBadRequest)
	}

	url := r.FormValue("url")

	hashedUrl, err := hasher.Sha256(url, hashLen)
	if err != nil {
		fmt.Println(err)
		return
	}

	shortUrl := defaultUrl + hashedUrl

	urlPair := storage.URLPair{
		LongUrl:  url,
		ShortUrl: shortUrl,
	}

	data, err := s.repo.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	data = append(data, urlPair)

	s.repo.Store(data)

	w.Write([]byte(shortUrl))
}
