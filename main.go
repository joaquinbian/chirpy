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

	mux.HandleFunc("GET /api/healthz", handleHealth)
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirpy)

	mux.Handle("/app/", cfg.middlewareIncHits(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", cfg.handleCountRequests)
	mux.HandleFunc("POST /admin/reset", cfg.handleResetCount)

	fmt.Println("server up and listening in 8080!")
	server.ListenAndServe()
}

func (cfg *apiConfig) handleCountRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/htlm; charset=utf-8")

	html := `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`

	w.Write([]byte(fmt.Sprintf(html, cfg.fileserverHits.Load())))
}
func (cfg *apiConfig) middlewareIncHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
