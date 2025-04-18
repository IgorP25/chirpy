package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/IgorP25/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find API key", err)
		return
	}

	if apiKey != cfg.polkaAPIKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeUserByID(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Could not retrieve user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Could not update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
