package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	params := parameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("error decoding params: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error processing body"})
		return
	}

	expirationTime := getExpTime(params.ExpiresInSeconds)

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		log.Printf("error getting user: %v", err)
		writeJSON(w, http.StatusUnauthorized, Envelope{"error": "invalid mail or password"})
		return
	}

	ok, err := auth.ComparePasswordHash(params.Password, user.HashedPassword)

	if !ok || err != nil {
		log.Printf("error getting user: %v", err)
		writeJSON(w, http.StatusUnauthorized, Envelope{"error": "invalid mail or password"})
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expirationTime)

	if err != nil {
		log.Printf("error creating token: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error creating token"})
		return
	}
	type response struct {
		User
		Token string `json:"token"`
	}

	writeJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})
}

func getExpTime(exp int) time.Duration {
	duration := (time.Duration(exp) * time.Second)

	if duration.Hours() == 0 || duration.Hours() > 1 {
		return time.Hour
	}

	return duration
}
