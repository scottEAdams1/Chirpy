package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/scottEAdams1/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		deleteDatabse()
	}
	//Create database
	db, err := database.NewDB("database.json")
	if err != nil {
		panic(err)
	}

	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
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
	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("POST /api/login", apiCfg.login)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.refresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.revoke)
	mux.HandleFunc("DELETE /api/chirps/{chirpid}", apiCfg.deleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.chirpyRed)

	//Create server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	//Run server
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func deleteDatabse() {
	// Define the path to the JSON file
	jsonFilePath := "./database.json"

	// Attempt to remove the file
	err := os.Remove(jsonFilePath)
	if err != nil {
		log.Fatal(err)
	}
}
