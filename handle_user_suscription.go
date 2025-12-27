package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleUserSuscription(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Event string
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("error decoding params: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error decoding params"})
		return
	}

	if params.Event != "user.upgraded" {
		writeJSON(w, http.StatusNoContent, Envelope{})
		return
	}

	_, err = cfg.db.GetUserByID(r.Context(), params.Data.UserID)

	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, Envelope{"error": "user not found"})
		return
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "something went wrong"})
		return
	}

	err = cfg.db.UpgradeUserChirpyRed(r.Context(), params.Data.UserID)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "something went wrong"})
		return
	}

	writeJSON(w, http.StatusNoContent, Envelope{})
}
