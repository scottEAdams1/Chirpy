package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	//Handlers
	mux := http.NewServeMux()
	fileServer := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/app/*", fileServer)
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	mux.HandleFunc("GET /api/reset", apiCfg.reset)
	mux.HandleFunc("POST /api/validate_chirp", validate)

	//Create server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	//Run server
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func validate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		msg := string(err.Error())
		respondWithError(w, 500, msg)
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	type success struct {
		Cleaned string `json:"cleaned_body"`
	}

	cleaned := getCleanedBody(params.Body)
	respondWithJSON(w, http.StatusOK, success{
		Cleaned: cleaned,
	})
}

func getCleanedBody(body string) string {
	words := strings.Split(body, " ")
	cleaned_words := []string{}
	for _, word := range words {
		lowerword := strings.ToLower(word)
		if lowerword == "kerfuffle" || lowerword == "sharbert" || lowerword == "fornax" {
			word = "****"
		}
		cleaned_words = append(cleaned_words, word)
	}
	cleaned := strings.Join(cleaned_words, " ")
	return cleaned
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code == 500 {
		log.Printf("Error decoding parameters: %s", msg)
	}
	type errorStruct struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorStruct{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	jsonResponse, jsonErr := json.Marshal(payload)
	if jsonErr != nil {
		log.Printf("Error marshalling JSON: %s", jsonErr)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(jsonResponse)
}
