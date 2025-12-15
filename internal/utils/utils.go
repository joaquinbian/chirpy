package utils

import (
	"slices"
	"strings"
)

var ProfaneWords []string = []string{"kerfuffle", "sharbert", "fornax"}

func IsProfane(word string) bool {

	loweredWord := strings.ToLower(word)

	return slices.Contains(ProfaneWords, loweredWord)

}

func CleanMessageProfane(message string) string {
	words := strings.Split(message, " ")

	for idx, w := range words {
		isProfane := IsProfane(w)
		if isProfane {
			words[idx] = "****"
		}

	}

	message = strings.Join(words, " ")

	return message
}
