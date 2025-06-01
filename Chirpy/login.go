package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Yssengrim/Chirpy/internal/auth"
	"github.com/Yssengrim/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}
	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the input parameters
	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
		return
	}
	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
		return
	}
	// Decide how long the token should last
	expiresIn := 3600
	if params.ExpiresInSeconds != nil && *params.ExpiresInSeconds >= 1 && *params.ExpiresInSeconds <= 3600 {
		expiresIn = *params.ExpiresInSeconds
	}

	// Generate the JWT token using the correct expiresIn value
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(expiresIn)*time.Second)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	// Generate a refresh token
	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}
	paramsForStore := database.StoreRefreshTokenParams{
		Token:     refresh_token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	}

	// Store the refresh token in the database
	if err := cfg.dbQueries.StoreRefreshToken(r.Context(), paramsForStore); err != nil {
		http.Error(w, "Failed to store refresh token", http.StatusInternalServerError)
		return
	}

	// Respond with the user info and the token
	response := struct {
		ID           string `json:"id"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
	}{
		ID:           user.ID.String(),
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
		Email:        user.Email,
		Token:        token,
		RefreshToken: refresh_token,
		IsChirpyRed:  user.IsChirpyRed,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
