package main

import (
	"chirpy/internal/auth"
	"database/sql"
	"log"
	"net/http"
	"time"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
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

	if time.Now().Compare(rToken.ExpiredAt) > 0 || err == sql.ErrNoRows || rToken.RevokedAt.Valid {
		log.Printf("error unauthorized token: %v", err)
		writeJSON(w, http.StatusUnauthorized, Envelope{"error": "invalid refresh token"})
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), rToken.Token)

	if err != nil {
		log.Printf("error getting refreshToken: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "something went wrong"})
		return
	}

	aToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, EXP_TIME_SECONDS)

	if err != nil {
		log.Printf("error creating JWT on refresh: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "something went wrong"})
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	if err != nil {
		log.Printf("error refresh and creating jwt: %v", err)
		writeJSON(w, http.StatusInternalServerError, Envelope{"error": "something went wrong"})
		return
	}

	writeJSON(w, http.StatusOK, response{
		Token: aToken,
	})
}
