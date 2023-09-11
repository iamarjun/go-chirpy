package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iamarjun/go-chirpy/internal/database"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request, db *database.DB) {
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

	userID, err := strconv.Atoi(id)

	if err != nil {
		respondWithError(w, 401, fmt.Sprintf(" %v", err))
		return
	}

	chirpId, err := strconv.Atoi(chi.URLParam(r, "chirpId"))

	if err != nil {
		respondWithJson(w, 400, fmt.Errorf("cannot convert %v to integer", chirpId))
		return
	}

	isDeleted, err := db.DeleteChirpsById(chirpId, userID)

	if err != nil {
		respondWithJson(w, 403, err)
		return
	}

	if !isDeleted {
		respondWithJson(w, 400, "chirp not found")
		return
	}

	respondWithJson(w, 200, http.StatusOK)
}
