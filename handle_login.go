package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("error decoding params: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error processing body"})
		return
	}

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

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, EXP_TIME_SECONDS)

	if err != nil {
		log.Printf("error creating access token: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "something went wrong"})
		return
	}

	rTokenStr, err := auth.MakeRefreshToken()

	if err != nil {
		log.Printf("error creating refresh token: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "something went wrong"})
		return
	}
	createRTokenParams := database.CreateRefreshTokenParams{
		Token:     rTokenStr,
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		ExpiredAt: time.Now().Add(time.Hour * 24 * 60),
	}

	rToken, err := cfg.db.CreateRefreshToken(r.Context(), createRTokenParams)

	if err != nil {
		log.Printf("error saving refresh token: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "something went wrong"})
		return
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	writeJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        token,
		RefreshToken: rToken.Token,
	})
}
