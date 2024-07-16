package main

import (
	"net/http"
	"strconv"
)

// Get all the chirps in the database
func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

// Get a single chirp from the database
func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	//Get chirpid from the path, then convert it to int
	chirpid := r.PathValue("chirpid")
	chirpidint, err := strconv.Atoi(chirpid)
	if err != nil {
		respondWithError(w, 400, err.Error())
	}

	//Get all chirps from the database
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//If the chirpid doesn't exist, return error
	if chirpidint > len(chirps) {
		respondWithError(w, 404, "Chirp doesn't exist")
		return
	}

	//Get chirp with chirpid
	chirp := chirps[chirpidint-1]
	respondWithJSON(w, http.StatusOK, chirp)
}
