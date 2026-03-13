package main

import (
	"fmt"
	hasher "urlShortener/pkg/hasher"
	storage "urlShortener/pkg/storage"
)

const (
	hashLen        = 10
	defaultUrl     = "https://trash-url.com/"
	defaultStorage = "urlShortenerDB.json"
)

func main() {
func homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./templates/index.html")
}
func shortenUrl(w http.ResponseWriter, r *http.Request) {
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

	urlRepo := storage.Repository{DbPath: defaultStorage}

	data, err := urlRepo.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	data = append(data, urlPair)

	urlRepo.Store(data)

	w.Write([]byte(shortUrl))
}
