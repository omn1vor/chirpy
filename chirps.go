package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/omn1vor/chirpy/internal/dto"
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
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while getting chirp: "+err.Error())
		return
	}
	if chirp == nil {
		respondWithError(w, http.StatusNotFound, "ID not found")
		return
	}
	respondWithJson(w, http.StatusOK, chirp)
}
