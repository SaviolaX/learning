package main

import (
	server "urlShortenerV1/internal/server"
	jsonRepo "urlShortenerV1/internal/repositorie/json"
)

const dbPath = "UrlShortenerDB.json"
const indexPath = "templates/index.html"

func main() {
	db := &jsonRepo.Repository{DbPath: dbPath}
	s := server.New(db, indexPath)
	s.Start(":8080")
}
	