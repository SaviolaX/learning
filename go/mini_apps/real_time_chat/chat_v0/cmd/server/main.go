package main

import (
	"chatV0/internal/handler"
	"chatV0/internal/hub"
	"log"
	"net/http"
)

var port = ":3000"
var host = "localhost"

func main() {
	h := hub.New()

	go h.Run()

	mux := http.NewServeMux()
	handler.New(h).RegisterRoutes(mux)

	log.Printf("server starts on: http://%s%s", host, port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
