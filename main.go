package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareIncHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {

	mux := http.NewServeMux()
	server := &http.Server{
		Handler:     mux,
		Addr:        ":8080",
		IdleTimeout: 10 * time.Second,
	}
	var cfg = &apiConfig{
		fileserverHits: atomic.Int32{},
	}
	mux.Handle("/app/", cfg.middlewareIncHits(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/metrics", cfg.handleCountRequests)
	mux.HandleFunc("/reset", cfg.handleResetCount)

	fmt.Println("server up and listening in 8080!")
	server.ListenAndServe()
}

func (cfg *apiConfig) handleCountRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handleResetCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	cfg.fileserverHits.Store(0)
	w.Write([]byte(fmt.Sprint("count reseted")))

}
