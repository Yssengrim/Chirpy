package main

import (
	"encoding/json"
	"net/http"
)

func (a *apiConfig) handlerAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := a.dbQueries.AllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	response := make([]chirpResponse, len(chirps))
	for i, chirp := range chirps {
		response[i] = chirpResponse{
			ID:        chirp.ID,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
