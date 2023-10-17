package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/omn1vor/chirpy/internal/auth"
	"github.com/omn1vor/chirpy/internal/dto"
)

func (cfg *apiConfig) addUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	userRequest := dto.UserRequest{}
	err := decoder.Decode(&userRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't decode user: "+err.Error())
		return
	}

	if strings.TrimSpace(userRequest.Email) == "" {
		respondWithError(w, http.StatusBadRequest, "Email cannot be empty")
		return
	}

	if userRequest.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password cannot be empty")
		return
	}

	if user, err := cfg.db.FindUserByEmail(userRequest.Email); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while working with DB: "+err.Error())
		return
	} else if user != nil {
		respondWithError(w, http.StatusBadRequest, "Email is already registered: "+userRequest.Email)
		return
	}

	hash, err := auth.HashPassword(userRequest.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash the password: "+err.Error())
		return
	}
	userRequest.Password = hash

	user, err := cfg.db.CreateUser(userRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not write user to the database: "+err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, user)
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	userRequest := dto.UserRequest{}
	err := decoder.Decode(&userRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't decode user: "+err.Error())
		return
	}

	if strings.TrimSpace(userRequest.Email) == "" {
		respondWithError(w, http.StatusBadRequest, "Email cannot be empty")
		return
	}

	user, err := cfg.db.FindUserByEmail(userRequest.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while working with DB: "+err.Error())
		return
	}
	if user == nil {
		respondWithError(w, http.StatusNotFound, "No users found with this email: "+userRequest.Email)
		return
	}
	if err = auth.CheckPassword(userRequest.Password, user.PwdHash); err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, user.ToDto())
}
