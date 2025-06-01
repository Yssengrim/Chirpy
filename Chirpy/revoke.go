package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return
	}

	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")
	if refreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusUnauthorized)
		return
	}

	refreshTokenDB, err := cfg.dbQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if refreshTokenDB.RevokedAt.Valid {
		http.Error(w, "Refresh token has already been revoked", http.StatusBadRequest)
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), refreshTokenDB.Token)
	if err != nil {
		http.Error(w, "Failed to revoke refresh token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
