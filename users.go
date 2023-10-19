package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/omn1vor/chirpy/internal/auth"
	"github.com/omn1vor/chirpy/internal/dto"
	"github.com/omn1vor/chirpy/internal/tokens"
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

	loginRequest := dto.LoginRequest{}
	err := decoder.Decode(&loginRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't decode user: "+err.Error())
		return
	}

	if strings.TrimSpace(loginRequest.Email) == "" {
		respondWithError(w, http.StatusBadRequest, "Email cannot be empty")
		return
	}

	user, err := cfg.db.FindUserByEmail(loginRequest.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while working with DB: "+err.Error())
		return
	}
	if user == nil {
		respondWithError(w, http.StatusNotFound, "No users found with this email: "+loginRequest.Email)
		return
	}
	if err = auth.CheckPassword(loginRequest.Password, user.PwdHash); err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userID := strconv.Itoa(user.Id)
	token, err := tokens.CreateAccessToken(cfg.jwtSecret, cfg.serviceId+"-access", userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while creating access token: "+err.Error())
		return
	}

	refreshToken, err := tokens.CreateRefreshToken(cfg.jwtSecret, cfg.serviceId+"-refresh", userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while creating refresh token: "+err.Error())
		return
	}

	loggedUserDto := struct {
		Id           string `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		Id:           userID,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}

	respondWithJson(w, http.StatusOK, loggedUserDto)
}

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Expecting JWT authorization")
		return
	}

	issuer := cfg.serviceId + "-access"
	userID, err := tokens.GetUserIdFromToken(cfg.jwtSecret, issuer, auth.GetTokenStringFromAuthHeader(authHeader))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong user ID format: "+userID)
		return
	}

	decoder := json.NewDecoder(r.Body)

	userRequest := dto.UserRequest{}
	err = decoder.Decode(&userRequest)
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

	hash, err := auth.HashPassword(userRequest.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash the password: "+err.Error())
		return
	}
	userRequest.Password = hash

	user, err := cfg.db.UpdateUser(id, userRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user: "+err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, user)
}

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Expecting JWT authorization")
		return
	}

	tokenString := auth.GetTokenStringFromAuthHeader(authHeader)
	issuer := cfg.serviceId + "-refresh"
	userID, err := tokens.GetUserIdFromToken(cfg.jwtSecret, issuer, tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if ok, err := cfg.db.TokenIsNotRevoked(tokenString); !ok {
		respondWithError(w, http.StatusUnauthorized, "Token is revoked")
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while checking token status")
		return
	}

	token, err := tokens.CreateAccessToken(cfg.jwtSecret, cfg.serviceId+"-access", userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while creating access token: "+err.Error())
		return
	}
	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	respondWithJson(w, http.StatusOK, response)
}

func (cfg *apiConfig) revokeToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Expecting JWT authorization")
		return
	}

	tokenString := auth.GetTokenStringFromAuthHeader(authHeader)
	issuer := cfg.serviceId + "-refresh"
	if _, err := tokens.GetUserIdFromToken(cfg.jwtSecret, issuer, tokenString); err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err := cfg.db.RevokeToken(tokenString); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
