package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		w.WriteHeader(500)
		log.Printf("failed to marshal JSON response %v", data)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	fmt.Printf("Data to respond %v\n", data)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5xx error: %s", msg)
	}

	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errResponse{
		Error: msg,
	})
}
