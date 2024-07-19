package main

import (
	"net/http"

	"github.com/scottEAdams1/Chirpy/internal/auth"
)

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	//Get the token string from the header
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = cfg.DB.RemoveToken(tokenString)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
