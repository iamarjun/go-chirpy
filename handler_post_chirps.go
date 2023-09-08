package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iamarjun/go-chirpy/internal/database"
)

func handlerPostChirps(w http.ResponseWriter, r *http.Request, db *database.DB) {
	type parameter struct {
		Body string `json:"body"`
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

	chirp, err := db.CreateChirp(params.Body)

	fmt.Printf("DB write done %v\n", chirp)
	if err != nil {
		respondWithJson(w, 400, err)
		return
	}
	
	fmt.Printf("Trying to respond with created chirp %v\n", chirp)
	respondWithJson(w, 201, chirp)

}
