package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) addUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	dto := userDto{}
	err := decoder.Decode(&dto)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't decode user: "+err.Error())
		return
	}

	if strings.TrimSpace(dto.Email) == "" {
		respondWithError(w, http.StatusBadRequest, "Email cannot be empty")
		return
	}

	user, err := cfg.db.CreateUser(dto.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't write user to the database: "+err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, user)
}
