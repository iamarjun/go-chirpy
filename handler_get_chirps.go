package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi"
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

func handlerGetChirpById(w http.ResponseWriter, r *http.Request, db *database.DB) {
	chirpId, err := strconv.Atoi(chi.URLParam(r, "chirpId"))

	if err != nil {
		respondWithJson(w, 400, fmt.Errorf("cannot convert %v to integer", chirpId))
		return
	}

	chirps, err := db.GetChirps()

	if err != nil {
		respondWithJson(w, 400, err)
		return
	}

	for _, chirp := range chirps {
		if chirp.ID == chirpId {
			respondWithJson(w, 200, chirp)
			return
		}
	}

	respondWithJson(w, 404, "")
}
