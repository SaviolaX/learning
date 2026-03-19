package main

import (
	application "urlShortenerV2/internal/application/url"
	"urlShortenerV2/internal/infrastructure/databases/json"
	h "urlShortenerV2/internal/infrastructure/hasher"
	transp "urlShortenerV2/internal/transport/http"
)

const dbPath = "UrlShortenerDB.json"
const indexPath = "templates/index.html"
const maxCodeLen = 10

var port = ":8080"
var host = "localhost"

func main() {
	repo := json.NewRepository(dbPath)

	service := application.NewService(repo, h.Sha256, maxCodeLen)

	handler := transp.NewHandler(indexPath, service, host, port)
	server := transp.NewServer(handler)
	server.Start(port)
}
