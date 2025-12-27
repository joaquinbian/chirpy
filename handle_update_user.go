package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handleEditUser(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Printf("error getting token: %v", err)
		writeJSON(w, http.StatusUnauthorized, Envelope{"error": "error getting token"})
		return
	}

	userID, err := auth.ValidateToken(token, cfg.jwtSecret)

	if err != nil {
		log.Printf("error validating token: %v", err)
		writeJSON(w, http.StatusUnauthorized, Envelope{"error": "error validating token"})
		return
	}

	type parameters struct {
		Email    string
		Password string
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&params)

	if err != nil {
		log.Printf("error decoding params: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "something went wrong editing user"})
		return
	}

	passHash, err := auth.HashPassword(params.Password)

	if err != nil {
		log.Printf("error hashing password: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "something went wrong editing user"})
		return
	}

	u, err := cfg.db.EditUser(r.Context(), database.EditUserParams{
		Email:          params.Email,
		HashedPassword: passHash,
		ID:             userID,
	})

	if err != nil {
		log.Printf("error editing user: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "something went wrong editing user"})
		return
	}

	type response struct {
		User
	}

	user := User{
		ID:          u.ID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Email:       u.Email,
		IsChirpyRed: u.IsChirpyRed,
	}

	writeJSON(w, http.StatusOK, response{
		User: user,
	})
}
