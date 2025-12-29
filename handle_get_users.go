package main

import (
	"log"
	"net/http"

	"github.com/joaquinbian/chirpy/internal/auth"
)

func (cfg *apiConfig) handleGetUsers(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Printf("error getting token: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error getting token"})
		return
	}

	_, err = auth.ValidateToken(token, cfg.jwtSecret)

	if err != nil {
		log.Printf("error validating token: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error validating token"})
		return
	}

	users, err := cfg.db.GetUsers(r.Context())

	if err != nil {
		log.Printf("error getting users: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error getting users"})
		return
	}

	writeJSON(w, http.StatusOK, users)
}
