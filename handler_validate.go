package main

import (
	"chirpy/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type parameters struct {
	Body string `json:"body"`
}

func handleValidateChirpy(w http.ResponseWriter, r *http.Request) {

	params := parameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)

	if err != nil {

		writeJSON(w, http.StatusBadRequest, Envelope{"Message": "Something went wrong"})
		return
	}

	if len(params.Body) >= 140 {

		err = writeJSON(w, http.StatusBadRequest, Envelope{"message": "Chip is too long"})

		if err != nil {
			fmt.Errorf("Error marshalling error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	cleanedWord := utils.CleanMessageProfane(params.Body)

	err = writeJSON(w, http.StatusOK, Envelope{"cleaned_body": cleanedWord})

	if err != nil {
		fmt.Errorf("Error marshalling error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
