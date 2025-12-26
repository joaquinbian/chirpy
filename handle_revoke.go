package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"database/sql"
	"log"
	"net/http"
	"time"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Printf("error getting refresh token: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "something went wrong"})
		return
	}

	rToken, err := cfg.db.GetRefreshToken(r.Context(), token)

	if err != nil {
		log.Printf("error getting refresh token from DB: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "something went wrong"})
		return
	}

	revokeTokenParams := database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Token: rToken.Token,
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), revokeTokenParams)

	if err != nil {
		log.Printf("error revoking token: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "something went wrong"})
		return

	}

	writeJSON(w, http.StatusNoContent, Envelope{})
}
