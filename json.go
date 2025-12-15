package main

import (
	"encoding/json"
	"net/http"
)

// Envelope: map donde el key es un string y el value es un interface vacio, que quiere decir que puede ser cualquier tipo de dato
type Envelope map[string]interface{}

func writeJSON(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	js = append(js, '\n') // Agrega un salto de línea al final del JSON

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}
