package handler

import (
	"chatV0/internal/client"
	"chatV0/internal/hub"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub *hub.Hub
}

func New(h *hub.Hub) *Handler {
	return &Handler{hub: h}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.serveHome)

	mux.HandleFunc("GET /ws", h.serveWs)
}

func (h *Handler) serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, "web/index.html")
}

func (h *Handler) serveWs(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name is required (?name=Alice)", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error for [%s]: %v", name, err)
		return
	}

	client.NewClient(conn, h.hub, name)
}
