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

func (cfg *apiConfig) handlerPutUsers(w http.ResponseWriter, r *http.Request, db *database.DB) {
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

	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type errorResp struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}


	err = decoder.Decode(&params)


	if err != nil {
		respondWithError(w, 500, fmt.Sprintf(" %v", err))
		return
	}

	if len(params.Email) == 0 {
		err := errorResp{}
		err.Error = "email cannot be empty"
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	if len(params.Password) == 0 {
		err := errorResp{}
		err.Error = "password cannot be empty"
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}


	userId, err := strconv.Atoi(id)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	if len(params.Email) == 0 {
		err := errorResp{}
		err.Error = "Email cannot be empty"
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	if len(params.Password) == 0 {
		err := errorResp{}
		err.Error = "Password cannot be empty"
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	isUpdated, user, err := db.UpdateUser(userId, params.Email, params.Password)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	if !isUpdated {
		respondWithError(w, 400, "something went wrongz")
		return
	}

	if err != nil {

		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	respondWithJson(w, 200, database.UserToResponseUser(user))

}
