package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iamarjun/go-chirpy/internal/database"
)

func handlerPostUsers(w http.ResponseWriter, r *http.Request, db *database.DB) {
	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type errorResp struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf(" %v", err))
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

	user, err := db.CreateUserWithPassword(params.Email, params.Password)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	respondWithJson(w, 201, database.UserToResponseUser(user))

}
