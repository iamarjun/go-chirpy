package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Body string `json:"body"`
	}

	type errorResp struct {
		Error string `json:"error"`
	}

	type successResp struct {
		CleanedBody string `json:"cleaned_body"`
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

	fmt.Println(len(params.Body))

	words := strings.Split(params.Body, " ")

	for i, word := range words {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax" {
			words[i] = "****"
		}
	}

	fmt.Println(words)

	body := strings.Join(words, " ")

	fmt.Println(body)

	resp := successResp{
		CleanedBody: body,
	}

	fmt.Println(resp)
	respondWithJson(w, 200, resp)

}
