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
	url := "https://github.com"
	// badUrl := ""

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

	fmt.Println("Hashed URL:", urlPair)

}
