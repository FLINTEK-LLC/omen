package api

import (
	"errors"
	"net/http"
)

func (s *Server) listGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := s.store.ListGroups(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, groups)
}

func (s *Server) getGroup(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	g, err := s.store.GetGroup(r.Context(), name)
	if err != nil {
		writeError(w, http.StatusNotFound, errors.New("group not found"))
		return
	}
	writeJSON(w, http.StatusOK, g)
}
