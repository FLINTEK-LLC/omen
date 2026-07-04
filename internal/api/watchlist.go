package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/FLINTEK-LLC/omen/internal/store"
)

func (s *Server) listWatchlist(w http.ResponseWriter, r *http.Request) {
	entries, err := s.store.ListWatchlist(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

type createWatchlistRequest struct {
	Pattern   string `json:"pattern"`
	Label     string `json:"label"`
	NotifyVia string `json:"notify_via"`
}

func (s *Server) createWatchlistEntry(w http.ResponseWriter, r *http.Request) {
	var req createWatchlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if req.Pattern == "" {
		writeError(w, http.StatusBadRequest, errors.New("pattern is required"))
		return
	}

	id, err := s.store.InsertWatchlistEntry(r.Context(), store.WatchlistEntry{
		Pattern:   req.Pattern,
		Label:     req.Label,
		NotifyVia: req.NotifyVia,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (s *Server) deleteWatchlistEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid id"))
		return
	}

	if err := s.store.DeleteWatchlistEntry(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
