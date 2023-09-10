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

func (cfg *apiConfig) handlerPostLogin(w http.ResponseWriter, r *http.Request, db *database.DB) {
	type parameter struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
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

	fmt.Println("Before trying to write data to db")

	user, err := db.GetUserByEmail(params.Email)

	if err != nil {
		fmt.Printf("Login user with password error %v", err)
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	isValidPassword, err := db.ValidatePasswordForUser(user, params.Password)

	if !isValidPassword {
		respondWithError(w, 401, fmt.Sprintf(" %v", err))
		return
	}

	//Issue a JWT token
	expire := params.ExpiresInSeconds

	if expire == 0 {
		expire = 24
	}

	expiresAt := time.Now().UTC().Add(time.Duration(expire) * time.Second)
	registerClaims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Subject:   strconv.Itoa(user.ID),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
	cfg.jwtClaims = registerClaims
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, registerClaims)
	fmt.Println("Signing token with secret")
	token, err := jwtToken.SignedString(cfg.jwtSecret)

	if err != nil {
		fmt.Printf("token %v", token)
		respondWithError(w, 400, fmt.Sprintf(" %v", err))
		return
	}

	usr := database.UserToResponseUserWithToken(user)
	usr.Token = token

	fmt.Printf("jwt token %v\n", token)
	fmt.Printf("Trying to respond with created user %v\n", usr)
	respondWithJson(w, 200, usr)
}
