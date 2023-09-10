package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iamarjun/go-chirpy/internal/database"
)

const ACCESS_ISSUER = "chirpy-access"
const REFRESH_ISSUER = "chirpy-refresh"

func (cfg *apiConfig) handlerPostLogin(w http.ResponseWriter, r *http.Request, db *database.DB) {
	type parameter struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
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

	user, err := db.GetUserByEmail(params.Email)

	fmt.Printf("user gotten from email %v", user)

	if err != nil {
		fmt.Printf("Login user with password error %v", err)
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	isValidPassword, err := db.ValidatePasswordForUser(user, params.Password)

	if !isValidPassword {
		fmt.Println("Password validation failed")
		respondWithError(w, 401, fmt.Sprintf(" %v", err))
		return
	}

	//Issue a JWT token

	accessRegisterClaims := jwt.RegisteredClaims{
		Issuer:    ACCESS_ISSUER,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Subject:   strconv.Itoa(user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
	}
	refreshRegisterClaims := jwt.RegisteredClaims{
		Issuer:    REFRESH_ISSUER,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Subject:   strconv.Itoa(user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * 24 * 60)),
	}

	cfg.accessJwtClaims = accessRegisterClaims
	cfg.refreshJwtClaims = refreshRegisterClaims

	accessJwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessRegisterClaims)
	fmt.Println("Signing token with secret")
	accessToken, err := accessJwtToken.SignedString(cfg.jwtSecret)

	if err != nil {
		fmt.Printf("access token %v", accessToken)
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	refreshJwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshRegisterClaims)
	fmt.Println("Signing token with secret")
	refreshToken, err := refreshJwtToken.SignedString(cfg.jwtSecret)

	if err != nil {
		fmt.Printf("refresh token %v", refreshToken)
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	usr := database.UserToResponseUserWithToken(user)
	usr.Token = accessToken
	usr.RefreshToken = refreshToken

	fmt.Printf("Trying to respond with created user %v\n", usr)

	respondWithJson(w, 200, usr)
}
