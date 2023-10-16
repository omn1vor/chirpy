package main

import (
	"encoding/json"
	"net/http"
)

const chirpMaxLen = 140

func (cfg *apiConfig) addChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	dto := chirpDto{}
	err := decoder.Decode(&dto)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't decode chirp body: "+err.Error())
		return
	}

	if len(dto.Body) > chirpMaxLen {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirp, err := cfg.db.CreateChirp(replaceProfanity(dto.Body))
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
