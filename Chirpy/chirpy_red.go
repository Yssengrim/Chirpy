package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Yssengrim/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type polkaWebhookPayload struct {
	Event string           `json:"event"`
	Data  polkaWebhookData `json:"data"`
}
type polkaWebhookData struct {
	UserID string `json:"user_id"`
}

func (apiConfig *apiConfig) handlerChirpyRed(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	if apiKey != apiConfig.polkaKey {
		http.Error(w, "Forbidden: Polka key cannot access this endpoint", http.StatusUnauthorized)
		return
	}

	var payload polkaWebhookPayload
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if payload.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if payload.Data.UserID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(payload.Data.UserID)
	if err != nil {
		http.Error(w, "Invalid User ID format", http.StatusBadRequest)
		return
	}
	_, err = apiConfig.dbQueries.ChirpyRedUpdate(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
