package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handleResetCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	cfg.fileserverHits.Store(0)
	w.Write([]byte(fmt.Sprint("count reseted")))

}
