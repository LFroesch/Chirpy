package main

import (
	"encoding/json"
	"net/http"

	"github.com/LFroesch/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type webhookRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	APIKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key format")
		return
	}
	if APIKey != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := webhookRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if params.Event == "user.upgraded" {
		userID, err := uuid.Parse(params.Data.UserID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}
		err = cfg.db.EnableChirpyRedByID(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
