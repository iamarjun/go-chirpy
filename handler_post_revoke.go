package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iamarjun/go-chirpy/internal/database"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request, db *database.DB) {
	authorizationToken := r.Header.Get("Authorization")

	if len(authorizationToken) == 0 {
		respondWithError(w, 400, fmt.Sprintln("jwt token not found"))
		return
	}

	splitAuth := strings.Split(authorizationToken, " ")

	if len(splitAuth) == 1 {
		respondWithError(w, 400, fmt.Sprintln("invalid jwt token"))
		return
	}

	jwtToken := splitAuth[1]

	token, err := jwt.ParseWithClaims(jwtToken, &cfg.refreshJwtClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	}, jwt.WithLeeway(2*time.Second))

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	if !token.Valid {
		respondWithError(w, 400, "invalid token")
		return
	}

	issuer, err := token.Claims.GetIssuer()

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	if issuer != REFRESH_ISSUER {
		respondWithError(w, 400, "invalid token")
		return
	}

	_, err = token.Claims.GetSubject()

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	isRevoked, err := db.RevokeRefreshToken(jwtToken)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	if !isRevoked {
		respondWithError(w, 400, "unable to revoke token")
		return
	}

	respondWithJson(w, 200, http.StatusOK)
}
