package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iamarjun/go-chirpy/internal/database"
)

func handlerPostLogin(w http.ResponseWriter, r *http.Request, db *database.DB) {
	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type errorResp struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}

	fmt.Println(params)

	err := decoder.Decode(&params)

	fmt.Println(params)

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

	fmt.Println("Before trying to write data to db")

	user, err := db.GetUserByEmail(params.Email)

	if err != nil {
		fmt.Printf("Login user with password error %v", err)
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	isValidPassword, err := db.ValidatePasswordForUser(user, params.Password)

	if !isValidPassword {
		respondWithError(w, 401, fmt.Sprintf(" %v", err))
		return
	}

	fmt.Printf("Trying to respond with created user %v\n", user)
	respondWithJson(w, 200, database.UserToResponseUser(user))
}
