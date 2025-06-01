package main

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/Yssengrim/Chirpy/internal/auth"
	"github.com/Yssengrim/Chirpy/internal/database"
	"github.com/google/uuid"
)

// Removed local DeleteChirpParams struct; use database.DeleteChirpParams instead.

func (apiCfg *apiConfig) handlerChirpDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}
		if !strings.HasPrefix(tokenString, "Bearer ") {
			http.Error(w, "Authorization header must start with 'Bearer '", http.StatusUnauthorized)
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		if tokenString == "" {
			http.Error(w, "JWT token missing", http.StatusUnauthorized)
			return
		}
		userID, err := auth.ValidateJWT(tokenString, apiCfg.jwtSecret)
		if err != nil {
			http.Error(w, "Invalid JWT token", http.StatusUnauthorized)
			return
		}
		chirpID := r.PathValue("chirpID")
		if chirpID == "" {
			http.Error(w, "Chirp ID is required", http.StatusBadRequest)
			return
		}

		chirpUUID, err := uuid.Parse(chirpID)
		if err != nil {
			http.Error(w, "Invalid Chirp ID format", http.StatusBadRequest)
			return
		}

		chirp, err := apiCfg.dbQueries.GetChirpByID(r.Context(), chirpUUID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Chirp not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if chirp.UserID != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		params := database.DeleteChirpParams{
			ID:     chirpUUID,
			UserID: userID,
		}

		err = apiCfg.dbQueries.DeleteChirp(r.Context(), params)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent) // 204 No Content
	}

}
