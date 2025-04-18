package main

import (
	"net/http"
	"sort"

	"github.com/IgorP25/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Cannot parse Chirp ID.", err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Cannot retrieve Chirp.", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	author_id := r.URL.Query().Get("author_id")

	chirps, err := func(u string) ([]database.Chirp, error) {
		if u == "" {
			return cfg.db.GetChirps(r.Context())
		}

		userID, err := uuid.Parse(u)
		if err != nil {
			return []database.Chirp{}, err
		}
		return cfg.db.GetChirpsByUser(r.Context(), userID)
	}(author_id)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot retrieve Chirps.", err)
		return
	}

	outChirps := []Chirp{}

	for _, chirp := range chirps {
		outChirps = append(outChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	sort_order := r.URL.Query().Get("sort")
	if sort_order == "desc" {
		sort.Slice(outChirps, func(i, j int) bool {
			return outChirps[i].CreatedAt.After(outChirps[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, outChirps)
}
