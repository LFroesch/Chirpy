package main

import (
	"net/http"

	"github.com/LFroesch/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	chirpID := r.PathValue("chirpID")
	parsedID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Input Format")
		return
	}

	dbChirp, err := cfg.db.GetChirpsByID(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp Not Found")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Request lacks valid authentication token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwt_secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}
	if userID != dbChirp.UserID {
		respondWithError(w, http.StatusForbidden, "User is not the author of the chirp")
		return
	}
	err = cfg.db.DeleteChirpByID(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
