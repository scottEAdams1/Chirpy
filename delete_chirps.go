package main

import (
	"net/http"
	"strconv"

	"github.com/scottEAdams1/Chirpy/internal/auth"
)

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	//Get the token string from the header
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	//Get the user ID from the token string
	userID, err := auth.GetUserID(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	//Get chirpid from the path, then convert it to int
	chirpid := r.PathValue("chirpid")
	chirpidint, err := strconv.Atoi(chirpid)
	if err != nil {
		respondWithError(w, 400, err.Error())
	}

	//Get chirps from database
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	//Confirm userID is same as on the chirp
	if chirps[chirpidint-1].AuthorID != userID {
		respondWithError(w, 403, "unauthorised user")
		return
	}

	//Remove the chirp from the database
	err = cfg.DB.RemoveChirp(chirpidint)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	w.WriteHeader(204)
}
