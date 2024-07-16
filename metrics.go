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
	//Set response to HTML
	w.Header().Add("Content-Type", "text/html; charset=utf-8")

	//Set status code to 200
	w.WriteHeader(http.StatusOK)

	//Set hits to number in cfg
	hits := cfg.fileserverHits

	//Write hitsText to the response body
	hitsText := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", hits)
	w.Write([]byte(hitsText))
}
