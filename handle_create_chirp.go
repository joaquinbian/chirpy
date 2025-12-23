package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"chirpy/internal/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type parameters struct {
	Body string `json:"body"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	params := parameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("error decoding chirp body")
		writeJSON(w, http.StatusInternalServerError, Envelope{"Message": "Something went wrong"})
		return
	}
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Printf("error GetBearerToken: %v", err)
		writeJSON(w, http.StatusUnauthorized, Envelope{"error": "invalid token"})
		return
	}

	id, err := auth.ValidateToken(token, cfg.jwtSecret)

	if err != nil {
		log.Printf("error ValidateToken: %v", err)
		writeJSON(w, http.StatusUnauthorized, Envelope{"error": "invalid token"})
		return
	}

	err = isChirpValid(params.Body)

	if err != nil {
		log.Printf("invalid chirp: %v", err)
		writeJSON(w, http.StatusBadRequest, Envelope{"error": err})
		return
	}
	cleanedChirp := utils.CleanMessageProfane(params.Body)

	c, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: id,
	})

	if err != nil {
		log.Printf("error creating chirp: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error creating chirp"})
		return
	}
	type response struct {
		Chirp
	}

	chirp := Chirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	}

	writeJSON(w, http.StatusCreated, response{Chirp: chirp})
}

func isChirpValid(body string) error {

	if len(body) >= 140 {

		return errors.New("chirp is too long")
	}

	return nil
}
