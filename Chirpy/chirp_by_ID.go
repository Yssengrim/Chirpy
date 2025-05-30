package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (a *apiConfig) handlerChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirpIDParsed, err := uuid.Parse(chirpID)
	if err != nil {
		http.Error(w, "Invalid chirp ID", http.StatusBadRequest)
		return
	}
	chirp, err := a.dbQueries.GetChirpByID(r.Context(), chirpIDParsed)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Chirp not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chirp); err != nil {
		http.Error(w, "Failed to encode chirp", http.StatusInternalServerError)
		return
	}
}
