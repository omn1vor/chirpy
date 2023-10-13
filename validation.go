package main

import (
	"encoding/json"
	"net/http"
)

const chirpMaxLen = 140

func validateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)

	chirp := chirpDto{}
	err := decoder.Decode(&chirp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(errorDto{
			Error: "Can't decode chirp body: " + err.Error(),
		})
		return
	}

	if len(chirp.Body) > chirpMaxLen {
		w.WriteHeader(http.StatusBadRequest)
		encoder.Encode(errorDto{
			Error: "Chirp is too long",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder.Encode(
		struct {
			Valid bool `json:"valid"`
		}{
			Valid: true,
		})
}
