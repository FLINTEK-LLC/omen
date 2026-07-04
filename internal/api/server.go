package api

import (
	"io/fs"
	"net/http"

	"github.com/FLINTEK-LLC/omen/internal/store"
)

// Server holds the dependencies shared by all HTTP handlers.
type Server struct {
	store *store.Store
	hub   *Hub
	mux   *http.ServeMux
}

// NewServer builds the HTTP handler for OMEN's REST API, SSE stream, and
// embedded frontend. frontend is served for any path not matched by /api/*.
func NewServer(st *store.Store, hub *Hub, frontend fs.FS) *Server {
	s := &Server{store: st, hub: hub}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/victims", s.listVictims)
	mux.HandleFunc("GET /api/victims/{id}", s.getVictim)
	mux.HandleFunc("GET /api/groups", s.listGroups)
	mux.HandleFunc("GET /api/groups/{name}", s.getGroup)
	mux.HandleFunc("GET /api/watchlist", s.listWatchlist)
	mux.HandleFunc("POST /api/watchlist", s.createWatchlistEntry)
	mux.HandleFunc("DELETE /api/watchlist/{id}", s.deleteWatchlistEntry)
	mux.HandleFunc("GET /api/stream", s.stream)
	mux.HandleFunc("GET /api/stats", s.getStats)
	mux.Handle("/", http.FileServerFS(frontend))

	s.mux = mux
	return s
}

// Handler returns the top-level HTTP handler.
func (s *Server) Handler() http.Handler {
	return s.mux
}

// BroadcastNewVictim pushes a new-victim event to all connected SSE clients.
func (s *Server) BroadcastNewVictim(v store.Victim) {
	s.hub.BroadcastJSON(v)
}
