package main

import (
	"net/http"

	"github.com/omn1vor/chirpy/internal/auth"
	"github.com/omn1vor/chirpy/internal/tokens"
)

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
