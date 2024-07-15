package main

import (
	"log"
	"net/http"
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
	mux.HandleFunc("GET /healthz", healthzHandler)
	mux.HandleFunc("GET /metrics", apiCfg.hitsHandler)
	mux.HandleFunc("GET /reset", apiCfg.reset)

	//Create server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	//Run server
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
