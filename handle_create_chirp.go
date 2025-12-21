package main

import (
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
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
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

	err = isChirpValid(params.Body)

	if err != nil {
		log.Printf("invalid chirp: %w", err)
		writeJSON(w, http.StatusBadRequest, Envelope{"error": err})
		return
	}

	cleanedChirp := utils.CleanMessageProfane(params.Body)

	_, err = cfg.db.GetUserByID(r.Context(), params.UserId)

	if err != nil {
		message := "user with given id not exists"
		log.Printf(message)
		writeJSON(w, http.StatusBadRequest, Envelope{"error": message})
		return
	}

	c, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: params.UserId,
	})

	if err != nil {
		log.Printf("error creating chirp: %w", err)
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
