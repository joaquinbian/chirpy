package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy_access"
)

func HashPassword(password string) (string, error) {

	pass, err := argon2id.CreateHash(password, argon2id.DefaultParams)

	if err != nil {
		log.Printf("error hashing password: %v", err)
		return "", err
	}

	return pass, nil
}

func ComparePasswordHash(password, hash string) (bool, error) {
	equal, err := argon2id.ComparePasswordAndHash(password, hash)

	if err != nil {
		log.Printf("error comparing passwords: %v", err)
		return false, nil
	}

	return equal, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	tokenSigned, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Printf("error signing token: %v", err)
		return "", err
	}

	return tokenSigned, nil
}

func ValidateToken(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected algorithm: %v", t.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		log.Printf("error validating token: %v", err)
		return uuid.Nil, errors.New("error validating token")
	}

	id, err := token.Claims.GetSubject()

	if err != nil {
		log.Printf("error getting subject")
		return uuid.Nil, errors.New("error validating token")
	}

	issuer, err := token.Claims.GetIssuer()

	if err != nil {
		log.Printf("error getting issuer")
		return uuid.Nil, errors.New("error validating token")
	}

	if issuer != string(TokenTypeAccess) {
		log.Printf("invalid issuer")
		return uuid.Nil, errors.New("invalid issuer")
	}

	parsedID, err := uuid.Parse(id)

	if err != nil {
		log.Printf("error getting subject")
		return uuid.Nil, errors.New("error validating token")
	}

	return parsedID, nil
}

func GetBearerToken(headers http.Header) (string, error) {

	auth := headers.Get("Authorization")

	if len(auth) == 0 {
		return "", errors.New("error no auth header provided")
	}

	authSplitted := strings.Split(auth, " ")

	if len(authSplitted) < 2 || authSplitted[0] != "Bearer" {
		return "", errors.New("invalid auth header provided")
	}

	token := authSplitted[1]

	return token, nil
}
