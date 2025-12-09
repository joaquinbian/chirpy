package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {

	mux := http.NewServeMux()
	server := &http.Server{
		Handler:     mux,
		Addr:        ":8080",
		IdleTimeout: 10 * time.Second,
	}

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("OK"))
	})

	fmt.Println("server up and listening in 8080!")
	server.ListenAndServe()
}
