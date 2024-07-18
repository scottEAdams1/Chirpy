package main

import "net/http"

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	//Get the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, 500, "Authorization header not found")
		return
	}

	//Get the string for the refresh token
	tokenString := authHeader[len("Bearer "):]

	err := cfg.DB.RemoveToken(tokenString)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
