package main

import (
	"fmt"
	hasher "urlShortener/pkg/hasher"
	storage "urlShortener/pkg/storage"
)

const (
	hashLen    = 10
	defaultUrl = "https://trash-url.com/"
)

func main() {
	url := "https://google.com"
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

	fmt.Println("Hashed URL:", urlPair)

}
