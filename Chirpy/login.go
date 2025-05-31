package main

import (
	"encoding/json"
	"net/http"

	"github.com/Yssengrim/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	// Respond with the authenticated user's details
	w.Header().Set("Content-Type", "application/json")

	responseUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	json.NewEncoder(w).Encode(responseUser)
}
