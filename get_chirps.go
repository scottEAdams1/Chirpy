package main

import (
	"net/http"
	"strconv"

	"github.com/scottEAdams1/Chirpy/internal/database"
)

// Get all the chirps in the database
func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	//Get chirps from database
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//Optional, set chirps to be return in ascending or descending order
	order := r.URL.Query().Get("sort")
	if order == "desc" {
		chirpsTemp := make([]database.Chirp, 0, len(chirps))
		for i := len(chirps) - 1; i > -1; i-- {
			chirpsTemp = append(chirpsTemp, chirps[i])
		}
		chirps = chirpsTemp
	}

	//Optional, get id to return only their chirps
	idString := r.URL.Query().Get("author_id")
	if idString == "" {
		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	//Get chirps with that id
	chirpsWithID := make([]database.Chirp, 0, len(chirps))
	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	for _, chirp := range chirps {
		if chirp.AuthorID == id {
			chirpsWithID = append(chirpsWithID, chirp)
		}
	}

	respondWithJSON(w, http.StatusOK, chirpsWithID)
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
