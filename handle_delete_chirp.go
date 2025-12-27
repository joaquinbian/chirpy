package main

import (
	"chirpy/internal/auth"
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {

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

	idParam := r.PathValue("chirpID")

	if len(idParam) < 1 {
		writeJSON(w, http.StatusBadRequest, Envelope{"error": "chirp id param not provided"})
		return
	}

	chirpID, err := uuid.Parse(idParam)

	if err != nil {
		log.Printf("error parsing chirp id: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "invalid chirp id"})
		return
	}

	chirpToDelete, err := cfg.db.GetChirpByID(r.Context(), chirpID)

	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, Envelope{"error": "chirp not found"})
		return
	}

	if err != nil {
		log.Printf("error getting chirp to delete: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error getting chirp to delete"})
		return
	}

	if userID != chirpToDelete.UserID {
		writeJSON(w, http.StatusForbidden, Envelope{"error": "forbidden request"})
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), chirpID)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "could not delete chirp"})
		return
	}

	writeJSON(w, http.StatusNoContent, Envelope{"message": "chirp deleted successfully"})
}
