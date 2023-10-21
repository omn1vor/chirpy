package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/omn1vor/chirpy/internal/dto"
	"github.com/omn1vor/chirpy/internal/errs"
)

func (cfg *apiConfig) upgradeToChirpyRed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Expecting an API key")
		return
	}

	if strings.TrimPrefix(authHeader, "ApiKey ") != os.Getenv("POLKA_API_KEY") {
		respondWithError(w, http.StatusUnauthorized, "Wrong API key")
		return
	}

	decoder := json.NewDecoder(r.Body)

	polkaRequest := dto.PolkaRequest{}
	err := decoder.Decode(&polkaRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't decode polka request: "+err.Error())
		return
	}

	if polkaRequest.Event != "user.upgraded" {
		w.WriteHeader(http.StatusOK)
		return
	}

	err = cfg.db.SetChirpyRed(polkaRequest.Data.UserId, true)
	var errNotFound *errs.ErrNotFound
	if err != nil {
		switch {
		case errors.As(err, &errNotFound):
			respondWithError(w, http.StatusNotFound, errNotFound.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
