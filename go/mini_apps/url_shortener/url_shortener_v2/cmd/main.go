package main

import (
	"log"
	application "urlShortenerV2/internal/application/url"
	sqlite "urlShortenerV2/internal/infrastructure/databases/sqlite"
	h "urlShortenerV2/internal/infrastructure/hasher"
	transp "urlShortenerV2/internal/transport/http"
)

const jsonDbPath = "UrlShortenerDB.json"
const sqliteDbPath = "urls.db"
const indexPath = "templates/index.html"
const maxCodeLen = 10

var port = ":8080"
var host = "localhost"

func main() {
	repo, err := sqlite.NewRepository(sqliteDbPath)
	if err != nil {
		log.Fatal(err)
	}

	service := application.NewService(repo, h.Sha256, maxCodeLen)

	handler := transp.NewHandler(indexPath, service, host, port)
	server := transp.NewServer(handler)
	server.Start(port)
}
