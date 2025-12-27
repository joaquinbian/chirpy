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

type createUserParams struct {
	Email    string
	Password string
}
type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handleCreateUsers(w http.ResponseWriter, r *http.Request) {

	params := createUserParams{}

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	err := decoder.Decode(&params)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error decoding json received"})
		log.Printf("error creating user: %v", err)
		return
	}

	if params.Password == "" {
		writeJSON(w, http.StatusBadRequest, Envelope{"error": "password is required"})
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)

	if err != nil {
		log.Printf("error hashing password: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error creating user"})
		return
	}

	u, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error creating user"})
		log.Printf("error creating user: %v", err)
		return
	}

	type repsonse struct {
		User
	}

	user := User{
		ID:          u.ID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Email:       u.Email,
		IsChirpyRed: u.IsChirpyRed,
	}

	writeJSON(w, http.StatusCreated, repsonse{
		User: user,
	})
}
