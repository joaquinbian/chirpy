package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleGetChirpByID(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("chirpID")

	if len(id) < 1 {
		log.Printf("error: chirpID not provided")
		writeJSON(w, http.StatusBadRequest, Envelope{"error": "chirpID not provided"})
		return
	}

	chirpID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("error parsing id: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error parsing chirp id"})
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)

	if err == sql.ErrNoRows {
		log.Printf("chirp not found")
		writeJSON(w, http.StatusNotFound, Envelope{"error": "chirp not found"})
		return
	}

	if err != nil {
		log.Printf("error getting chirp: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error getting chirp"})
		return
	}

	writeJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.db.GetChirps(r.Context())

	if err != nil {
		log.Printf("error getting chirps: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error getting chirps"})
		return
	}

	var parsedChirps = []Chirp{}

	for _, c := range chirps {
		parsedChirps = append(parsedChirps, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}

	writeJSON(w, http.StatusOK, parsedChirps)
}
