package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iamarjun/go-chirpy/internal/database"
)

func handlerPostUsers(w http.ResponseWriter, r *http.Request, db *database.DB) {
	type parameter struct {
		Email string `json:"email"`
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
		respondWithJson(w, 500, err)
		return
	}

	if len(params.Email) == 0 {
		err := errorResp{}
		err.Error = "Email not found"
		respondWithJson(w, 400, err)
		return
	}

	fmt.Println("Before trying to write data to db")

	user, err := db.CreateUser(params.Email)

	fmt.Printf("DB write done %v\n", user)
	if err != nil {
		respondWithJson(w, 400, err)
		return
	}

	fmt.Printf("Trying to respond with created user %v\n", user)
	respondWithJson(w, 201, user)

}
