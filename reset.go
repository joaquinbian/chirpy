package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handleResetState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	cfg.fileserverHits.Store(0)

	err := cfg.db.DeleteUsers(r.Context())

	if err != nil {
		log.Printf("error deleting users: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "error deleting users"})
		return
	}

	writeJSON(w, http.StatusOK, Envelope{"message": "state reseted"})

}
