package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string
		Password string
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

	type response struct {
		User
	}

	writeJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
