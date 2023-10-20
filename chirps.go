package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/omn1vor/chirpy/internal/dto"
	"github.com/omn1vor/chirpy/internal/errs"
)

const chirpMaxLen = 140

func (cfg *apiConfig) addChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, err := cfg.getAuthenticatedUserID(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	authorID, err := strconv.Atoi(userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)

	chirpDto := dto.ChirpDto{}
	err = decoder.Decode(&chirpDto)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't decode chirp body: "+err.Error())
		return
	}

	if len(chirpDto.Body) > chirpMaxLen {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirpDto.AuthorId = authorID
	chirpDto.Body = replaceProfanity(chirpDto.Body)
	chirp, err := cfg.db.CreateChirp(chirpDto)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't write chirp to the database: "+err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	chirps, err := cfg.db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while getting chirps: "+err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Wrong chirp ID")
		return
	}

	chirp, err := cfg.db.GetChirp(id)
	var errNotFound *errs.ErrNotFound
	switch {
	case errors.As(err, &errNotFound):
		respondWithError(w, http.StatusNotFound, errNotFound.Error())
	default:
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJson(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	chirpID, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Wrong chirp ID: "+id)
		return
	}

	userID, err := cfg.getAuthenticatedUserID(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	authorID, err := strconv.Atoi(userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong author ID: "+userID+": "+err.Error())
		return
	}

	chirp, err := cfg.db.GetChirp(chirpID)
	var errNotFound *errs.ErrNotFound
	if err != nil {
		switch {
		case errors.As(err, &errNotFound):
			respondWithError(w, http.StatusNotFound, errNotFound.Error())
			return
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if chirp.AuthorId != authorID {
		respondWithError(w, http.StatusForbidden, "You can only delete your chirps")
		return
	}

	err = cfg.db.DeleteChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
