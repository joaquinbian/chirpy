package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.db.GetChirps(r.Context())

	if err != nil {
		log.Printf("error getting chirps: %w", err)
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
