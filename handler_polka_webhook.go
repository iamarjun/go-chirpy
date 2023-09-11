package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/iamarjun/go-chirpy/internal/database"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request, db *database.DB) {
	apiKey := r.Header.Get("Authorization")

	if len(apiKey) == 0 {
		respondWithError(w, 401, fmt.Sprintln("api key not found"))
		return
	}

	splitAuth := strings.Split(apiKey, " ")

	if len(splitAuth) == 1 {
		respondWithError(w, 401, fmt.Sprintln("invalid api key"))
		return
	}

	polkaApiKey := splitAuth[1]

	if polkaApiKey != cfg.polkaApiKey {
		respondWithError(w, 401, fmt.Sprintln("invalid api key"))
		return
	}

	type body struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := body{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	if strings.ToLower(params.Event) != "user.upgraded" {
		respondWithJson(w, 200, "")
		return
	}

	userId := params.Data.UserID
	isMarked, err := db.MarkUserAsChirpRed(userId)

	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}

	if !isMarked {
		respondWithError(w, 404, err.Error())
		return
	}

	respondWithJson(w, 200, "")

}
