package main

import (
	"log"
	"net/http"
	"sort"

	"github.com/iamarjun/go-chirpy/internal/database"
)

func handlerGetChirps(w http.ResponseWriter, r *http.Request, db *database.DB) {

	chirps, err := db.GetChirps()

	sort.SliceStable(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	if err != nil {
		log.Fatal(err)
		respondWithJson(w, 400, err)
		return
	}

	respondWithJson(w, 200, chirps)

}
