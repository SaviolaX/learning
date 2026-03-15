package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	hasher "urlShortenerV1/internal/hasher"
	jsonRepo "urlShortenerV1/internal/repositorie/json"
	url "urlShortenerV1/internal/url"
)

const maxCodeLen = 10
const dbPath = "UrlShortenerDB.json"
const redirectUrl = "http://localhost:8080/r/"
const indexPath = "templates/index.html"

type Server struct {
	urls      []jsonRepo.UrlPair
	db        jsonRepo.Repository
	indexPath string
}

func main() {

	db := jsonRepo.Repository{
		DbPath: dbPath,
	}

	loadedUrls, err := db.Load()
	if err != nil {
		log.Println(err)
	}

	s := &Server{urls: loadedUrls, db: db, indexPath: indexPath}

	router := http.NewServeMux()

	router.HandleFunc("/", s.homePage)
	router.HandleFunc("/shorten", s.ShortenUrl)
	router.HandleFunc("/r/", s.RedirectUrl)

	port := ":8080"

	fmt.Println("Starting server on port", port)

	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("server stopped: %v", err)
	}

}

func (s *Server) homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, s.indexPath)
}

func (s *Server) ShortenUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "cannot parse form", http.StatusBadRequest)
		return
	}

	origUrl := r.FormValue("url")

	validatedUrl, err := url.NewUrl(origUrl)
	if err != nil {
		log.Println(err)
		http.Error(w, "validation error", http.StatusInternalServerError)
		return
	}

	hashCode, err := hasher.Sha256(validatedUrl.Value, maxCodeLen)
	if err != nil {
		log.Println(err)
		http.Error(w, "hashing error", http.StatusInternalServerError)
		return
	}

	checkedUrl, ok := s.db.IsExist(hashCode, s.urls)
	// if not exist -> add, store, return
	if !ok {
		pair := jsonRepo.UrlPair{
			Code:        hashCode,
			OriginalUrl: origUrl,
		}

		s.urls = append(s.urls, pair)

		if err := s.db.Store(s.urls); err != nil {
			log.Println(err)
			http.Error(w, "storing error", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(redirectUrl + hashCode))
		return

	}
	w.Write([]byte(redirectUrl + checkedUrl.Code))
}

func (s *Server) RedirectUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hashCode := strings.TrimPrefix(r.URL.Path, "/r/")

	urlPair, ok := s.db.IsExist(hashCode, s.urls)
	fmt.Println(urlPair)

	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, urlPair.OriginalUrl, http.StatusFound)
}
