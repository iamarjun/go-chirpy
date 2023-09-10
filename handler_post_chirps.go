package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iamarjun/go-chirpy/internal/database"
)

func (cfg *apiConfig) handlerPostChirps(w http.ResponseWriter, r *http.Request, db *database.DB) {

	authorizationToken := r.Header.Get("Authorization")

	if len(authorizationToken) == 0 {
		respondWithError(w, 401, fmt.Sprintln("jwt token not found"))
		return
	}

	splitAuth := strings.Split(authorizationToken, " ")

	if len(splitAuth) == 1 {
		respondWithError(w, 401, fmt.Sprintln("invalid jwt token"))
		return
	}

	jwtToken := splitAuth[1]

	token, err := jwt.ParseWithClaims(jwtToken, &cfg.accessJwtClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	}, jwt.WithLeeway(2*time.Second))

	if err != nil {
		respondWithError(w, 401, fmt.Sprintf(" %v", err))
		return
	}

	if !token.Valid {
		respondWithError(w, 401, "invalid token")
		return
	}

	issuer, err := token.Claims.GetIssuer()

	if err != nil {
		respondWithError(w, 401, fmt.Sprintf(" %v", err))
		return
	}

	if issuer != ACCESS_ISSUER {
		respondWithError(w, 401, "invalid token")
		return
	}

	id, err := token.Claims.GetSubject()

	if err != nil {
		respondWithError(w, 401, fmt.Sprintf(" %v", err))
		return
	}

	userID, err := strconv.Atoi(id)

	if err != nil {
		respondWithError(w, 401, fmt.Sprintf(" %v", err))
		return
	}

	type parameter struct {
		Body string `json:"body"`
	}

	type errorResp struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}

	fmt.Println(params)

	err = decoder.Decode(&params)

	fmt.Println(params)

	if err != nil {
		respondWithJson(w, 500, err)
		return
	}

	if len(params.Body) == 0 {
		err := errorResp{}
		err.Error = "Chirp not found"
		respondWithJson(w, 400, err)
		return
	}

	if len(params.Body) > 144 {
		err := errorResp{}
		err.Error = "Chirp is too long"
		fmt.Println(err)
		respondWithJson(w, 400, err)
		return
	}

	fmt.Println("Before trying to write data to db")

	chirp, err := db.CreateChirp(params.Body, userID)

	fmt.Printf("DB write done for chirps endpoint %v\n", chirp)
	if err != nil {
		respondWithJson(w, 400, err)
		return
	}

	fmt.Printf("Trying to respond with created chirp %v\n", chirp)
	respondWithJson(w, 201, chirp)

}
