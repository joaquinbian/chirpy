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

	mux.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Println("server up and listening in 8080!")
	server.ListenAndServe()
}
