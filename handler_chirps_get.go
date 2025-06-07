package main

import (
	"net/http"

	"github.com/LFroesch/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	sortOrder := r.URL.Query().Get("sort")
	var dbChirps []database.Chirp
	var err error

	if sortOrder == "desc" {
		dbChirps, err = cfg.db.GetChirpsDesc(r.Context())
	} else {
		dbChirps, err = cfg.db.GetChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}
	chirps := []Chirp{}
	s := r.URL.Query().Get("author_id")
	if s != "" {
		for _, dbChirp := range dbChirps {
			if dbChirp.UserID.String() == s {
				chirps = append(chirps, Chirp{
					ID:        dbChirp.ID,
					CreatedAt: dbChirp.CreatedAt,
					UpdatedAt: dbChirp.UpdatedAt,
					UserID:    dbChirp.UserID,
					Body:      dbChirp.Body,
				})
			}
		}
	} else {
		for _, dbChirp := range dbChirps {
			chirps = append(chirps, Chirp{
				ID:        dbChirp.ID,
				CreatedAt: dbChirp.CreatedAt,
				UpdatedAt: dbChirp.UpdatedAt,
				UserID:    dbChirp.UserID,
				Body:      dbChirp.Body,
			})
		}
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirpsByID(w http.ResponseWriter, r *http.Request) {

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

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
