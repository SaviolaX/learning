package http

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	h *Handler
}

func NewServer(h *Handler) *Server {
	return &Server{
		h: h,
	}
}

func (s *Server) Start(port string) error {

	router := http.NewServeMux()

	router.HandleFunc("/", s.h.homePage)
	router.HandleFunc("/shorten", s.h.ShortenUrl)
	router.HandleFunc("/r/", s.h.RedirectUrl)

	fmt.Println("starting server on port", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("server stopped: %v", err)
	}

	return nil

}
