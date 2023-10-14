package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

const chirpMaxLen = 140

func validateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	chirp := chirpDto{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't decode chirp body: "+err.Error())
		return
	}

	if len(chirp.Body) > chirpMaxLen {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	const profaneMask = "****"
	badWords := profaneWords()
	allWords := strings.Split(chirp.Body, " ")
	for i, word := range allWords {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			allWords[i] = profaneMask
		}
	}

	encoder := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	encoder.Encode(
		struct {
			CleanedBody string `json:"cleaned_body"`
		}{
			CleanedBody: strings.Join(allWords, " "),
		})
}
