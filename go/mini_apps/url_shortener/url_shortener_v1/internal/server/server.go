package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	jsonRepo "urlShortenerV1/internal/repositorie/json"
	hasher "urlShortenerV1/internal/hasher"
	url "urlShortenerV1/internal/url"
)

const maxCodeLen = 10
const redirectUrl = "http://localhost:8080/r/"
const homePage = "http://localhost:8080/"


type Server struct {
	db        *jsonRepo.Repository
	indexPath string
}

func New(db *jsonRepo.Repository,indexPath string) *Server {
	return &Server{db: db, indexPath: indexPath}
}

func (s *Server) Start(port string) error {

	router := http.NewServeMux()

	router.HandleFunc("/", s.homePage)
	router.HandleFunc("/shorten", s.ShortenUrl)
	router.HandleFunc("/r/", s.RedirectUrl)

	fmt.Println("Starting server on port", port)
	fmt.Println("Home page: ", homePage)

	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("server stopped: %v", err)
	}
	return nil
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

	checkedUrl, err := s.db.GetByCode(hashCode)
	if err != nil && !errors.Is(err, jsonRepo.ErrNotFound) {
		log.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, jsonRepo.ErrNotFound) {

		log.Println("saving a new url pair")
		pair := jsonRepo.UrlPair{
			Code:        hashCode,
			OriginalUrl: validatedUrl.Value,
		}

		if err := s.db.Store(pair); err != nil {
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

	urlPair, err := s.db.GetByCode(hashCode)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, urlPair.OriginalUrl, http.StatusFound)
}
