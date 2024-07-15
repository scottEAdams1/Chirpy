package main

import (
	"fmt"
	"net/http"
)

// Increase number of server hits every time someone visits
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
	return handler
}

// Return number of hits to the website
func (cfg *apiConfig) hitsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileserverHits
	hitsText := fmt.Sprintf("Hits: %d", hits)
	w.Write([]byte(hitsText))
}
