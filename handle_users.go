package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type createUserParams struct {
	Email string
}
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handleCreateUsers(w http.ResponseWriter, r *http.Request) {

	params := createUserParams{}

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	err := decoder.Decode(&params)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, Envelope{"error": "error decoding json received"})
		log.Printf("error creating user: %v", err)
		return
	}

	u, err := cfg.db.CreateUser(r.Context(), params.Email)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, Envelope{"error": "error creating user"})
		log.Printf("error creating user: %v", err)
		return
	}

	user := User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}

	writeJSON(w, http.StatusCreated, Envelope{"user": user})
}
