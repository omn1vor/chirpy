package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const defaultTTL = 24 * 60 * 60

func CreateToken(sign, issuer, subject string, expiresInSeconds int) (string, error) {
	if expiresInSeconds == 0 {
		expiresInSeconds = defaultTTL
	}
	expiresInSeconds = min(expiresInSeconds, defaultTTL)

	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   subject,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expiresInSeconds))),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(sign))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func GetUserIdFromToken(sign, tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(sign), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	} else {
		return "", errors.New("invalid token")
	}
}
