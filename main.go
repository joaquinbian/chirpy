package main

import (
	"chirpy/internal/database"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	mode           string
}

func main() {
	godotenv.Load()

	dbUrl := os.Getenv("DB_URL")
	mode := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatalf("error connecting to db: %w", err)
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Handler:     mux,
		Addr:        ":8080",
		IdleTimeout: 10 * time.Second,
	}

	queries := database.New(db)

	var cfg = &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             queries,
		mode:           mode,
	}

	mux.HandleFunc("GET /api/healthz", handleHealth)
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirpy)
	mux.HandleFunc("POST /api/users", cfg.handleCreateUsers)

	mux.Handle("/app/", cfg.middlewareIncHits(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", cfg.handleCountRequests)
	mux.Handle("POST /admin/reset", cfg.middlewareDevMode(http.HandlerFunc(cfg.handleResetState)))

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

func (cfg *apiConfig) middlewareDevMode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cfg.mode != "dev" {
			writeJSON(w, http.StatusForbidden, Envelope{"error": "forbidden request"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
