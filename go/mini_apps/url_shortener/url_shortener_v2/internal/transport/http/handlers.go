package http

import (
	"net/http"
	"strings"
	application "urlShortenerV2/internal/application/url"
)

type Handler struct {
	indexPath string
	service   *application.Service
	host      string
	port      string
}

func NewHandler(indexPath string, service *application.Service, host string, port string) *Handler {
	return &Handler{
		indexPath: indexPath,
		service:   service,
		host:      host,
		port:      port,
	}
}

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, h.indexPath)
}

func (h *Handler) ShortenUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	origUrl := r.FormValue("url")

	code, err := h.service.Shorten(r.Context(), origUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	redirectUrl := "http://" + h.host + h.port + "/r/" + code

	w.Write([]byte(redirectUrl))
}

func (h *Handler) RedirectUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hashCode := strings.TrimPrefix(r.URL.Path, "/r/")

	url, err := h.service.GetUrl(r.Context(), hashCode)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)

}
