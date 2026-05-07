package server

import (
	"context"
	"encoding/json"
	"maps"
	"net/http"
	"slices"

	data "shinobi_showdown/internal/game/data"
)

type DataHandler struct {
	mux *http.ServeMux
}

func NewDataHandler(ctx context.Context) *DataHandler {
	handler := &DataHandler{
		mux: http.NewServeMux(),
	}

	return handler
}

func (dh *DataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dh.mux.ServeHTTP(w, r)
}

func (dh *DataHandler) HandleGetActors(w http.ResponseWriter, r *http.Request) {
	actors := slices.Collect(maps.Values(data.ACTORS))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(actors); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (dh *DataHandler) HandleGetActions(w http.ResponseWriter, r *http.Request) {
	actions := slices.Collect(maps.Values(data.ACTIONS))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(actions); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (dh *DataHandler) HandleGetItems(w http.ResponseWriter, r *http.Request) {
	actions := slices.Collect(maps.Values(data.ITEMS))
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(actions); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
