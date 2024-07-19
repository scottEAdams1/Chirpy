package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/scottEAdams1/Chirpy/internal/auth"
)

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	//Get the token string from the header
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	//Get the token from the database
	token, err := cfg.DB.GetToken(tokenString)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}
	if token.Expiration.After(time.Now()) == false {
		respondWithError(w, 401, "token expired")
		return
	}

	//Set the claims of the JWT
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		Subject:   strconv.Itoa(token.UserID),
	}

	//Create and sign the token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := accessToken.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	type Token struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, 200, Token{
		Token: signedToken,
	})
}
