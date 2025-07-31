package server

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"github.com/user/urlshortener/internal/core"
)

// Server is the HTTP server for the URL shortener.
type Server struct {
	core *core.Core
}

// NewServer creates a new Server instance.
func NewServer(core *core.Core) *Server {
	return &Server{core: core}
}

// RegisterRoutes registers the HTTP routes for the server.
func (s *Server) RegisterRoutes(router *httprouter.Router) {
	router.POST("/shorten", s.handleShorten)
	router.ServeFiles("/static/*filepath", http.Dir("web/static"))

	router.NotFound = http.HandlerFunc(s.handleRedirect)
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

func (s *Server) handleShorten(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortCode, err := s.core.ShortenURL(req.URL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to shorten URL")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

    baseURL := "http://" + r.Host + "/"
	res := shortenResponse{ShortURL: baseURL + shortCode}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error().Err(err).Msg("Failed to write response")
	}
}

func (s *Server) handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[1:]

	originalURL, err := s.core.GetURL(shortCode)
	if err != nil {
		// Handle not found error
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}
