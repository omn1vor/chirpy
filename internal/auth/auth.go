package auth

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash), err
}

func GetTokenStringFromAuthHeader(authHeader string) string {
	return strings.TrimPrefix(authHeader, "Bearer ")
}
