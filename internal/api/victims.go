package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/FLINTEK-LLC/omen/internal/store"
)

type victimDetail struct {
	store.Victim
	ShodanSnapshots []store.ShodanSnapshot `json:"shodan_snapshots"`
	KEVMatches      []store.KEVMatch       `json:"kev_matches"`
}

func (s *Server) listVictims(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := store.VictimFilter{
		Country: q.Get("country"),
		Group:   q.Get("group"),
		Since:   q.Get("since"),
	}
	if v, err := strconv.Atoi(q.Get("limit")); err == nil {
		filter.Limit = v
	}
	if v, err := strconv.Atoi(q.Get("offset")); err == nil {
		filter.Offset = v
	}

	victims, err := s.store.ListVictims(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, victims)
}

func (s *Server) getVictim(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	v, err := s.store.GetVictim(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, errors.New("victim not found"))
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	snapshots, err := s.store.ShodanSnapshotsForVictim(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	kevMatches, err := s.store.KEVMatchesForVictim(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, victimDetail{
		Victim:          v,
		ShodanSnapshots: snapshots,
		KEVMatches:      kevMatches,
	})
}

func (s *Server) getStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.store.GetStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}
