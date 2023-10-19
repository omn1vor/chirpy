package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const accessTokenExpiresIn = time.Hour
const refreshTokenExpiresIn = 60 * 24 * time.Hour

var errInvalidToken = errors.New("invalid token")

func CreateAccessToken(sign, issuer, subject string) (string, error) {
	return createToken(sign, issuer, subject, accessTokenExpiresIn)
}

func CreateRefreshToken(sign, issuer, subject string) (string, error) {
	return createToken(sign, issuer, subject, refreshTokenExpiresIn)
}

func createToken(sign, issuer, subject string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   subject,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(sign))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func GetUserIdFromToken(sign, issuer, tokenString string) (string, error) {
	token, err := getTokenFromString(sign, tokenString)
	if err != nil {
		return "", err
	}

	claims, err := getValidatedClaims(token, issuer)
	if err != nil {
		return "", nil
	}

	return claims.Subject, nil
}

func ValidateTokenString(sign, issuer, tokenString string) (bool, error) {
	token, err := getTokenFromString(sign, tokenString)
	if err != nil {
		return false, err
	}
	_, err = getValidatedClaims(token, issuer)
	if err != nil {
		return false, err
	}

	return true, nil
}

func getTokenFromString(sign, tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(sign), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func getValidatedClaims(token *jwt.Token, issuer string) (*jwt.RegisteredClaims, error) {
	if !token.Valid {
		return nil, errInvalidToken
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || claims.Issuer != issuer {
		return nil, errInvalidToken
	}

	return claims, nil
}
