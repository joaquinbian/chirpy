package auth

import (
	"log"

	"github.com/alexedwards/argon2id"
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
