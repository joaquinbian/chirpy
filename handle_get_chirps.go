package main

import (
	"chirpy/internal/database"
	"database/sql"
	"log"
	"net/http"
	"sort"

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

	queryID := r.URL.Query().Get("author_id")

	sort := "asc"
	sortParam := r.URL.Query().Get("sort")

	if sortParam == "desc" {
		sort = "desc"
	}

	if queryID != "" {

		userID, err := uuid.Parse(queryID)

		if err != nil {
			log.Printf("error parsing user id: %v", err)
			writeJSON(w, http.StatusBadRequest, Envelope{"error": "invalid author id"})
			return
		}

		chirps, err := cfg.db.GetChirpByUserID(r.Context(), userID)

		parsedChirps := parseChirps(chirps)
		sortedChirps := sortChirps(parsedChirps, sort)

		if err != nil {
			log.Printf("error getting chrips: %v", err)
			writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error getting chirps"})
			return
		}

		writeJSON(w, http.StatusOK, sortedChirps)
		return
	}

	chirps, err := cfg.db.GetChirps(r.Context())

	parsedChirps := parseChirps(chirps)
	sortedChirps := sortChirps(parsedChirps, sort)

	if err != nil {
		log.Printf("error getting chirps: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error getting chirps"})
		return
	}

	writeJSON(w, http.StatusOK, sortedChirps)
}

func parseChirps(dbChirps []database.Chirp) []Chirp {
	var parsedChirps = []Chirp{}

	for _, c := range dbChirps {
		parsedChirps = append(parsedChirps, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}

	return parsedChirps
}

func sortChirps(chirps []Chirp, mode string) []Chirp {
	if mode == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	} else if mode == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)

		})
	}

	return chirps
}
