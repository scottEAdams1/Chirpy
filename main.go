package main

import (
	"log"
	"net/http"

	"github.com/scottEAdams1/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	//Create database
	db, err := database.NewDB("database.json")
	if err != nil {
		panic(err)
	}

	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	//Handlers
	mux := http.NewServeMux()
	fileServer := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/app/*", fileServer)

	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	mux.HandleFunc("GET /api/reset", apiCfg.reset)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpid}", apiCfg.getChirp)

	//Create server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	//Run server
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
