package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Yssengrim/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
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
	if refreshTokenDB.ExpiresAt.Before(time.Now()) {
		http.Error(w, "Refresh token has expired", http.StatusUnauthorized)
		return
	}
	if refreshTokenDB.RevokedAt.Valid {
		http.Error(w, "Refresh token has been revoked", http.StatusUnauthorized)
		return
	}
	_, err = cfg.dbQueries.GetUserById(r.Context(), refreshTokenDB.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	accessToken, err := auth.MakeJWT(refreshTokenDB.UserID, cfg.jwtSecret, 3600*time.Second) // 1 hour
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	type response struct {
		Token string `json:"token"`
	}
	resp := response{
		Token: accessToken,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil { // Now this err is from the Encode call
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
