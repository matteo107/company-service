package main

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// createAuthenticationTokenHandler is a handler function which handles requests for creating a new token
func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// should go to DB and validate user, hardcoded for now
	if app.validateUser(input.Email, input.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Create JWT token
	IssuedAt := time.Now()
	ExpiresAt := IssuedAt.AddDate(0, 0, 1).Unix()
	NotBefore := IssuedAt.Unix()

	claims := &Claims{
		Username: input.Email,
		StandardClaims: jwt.StandardClaims{
			Audience:  "api.companyservice.io",
			ExpiresAt: ExpiresAt,
			Id:        "",
			IssuedAt:  IssuedAt.Unix(),
			Issuer:    "api.companyservice.io",
			NotBefore: NotBefore,
			Subject:   input.Email,
		},
	}

	// Create a new token object with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// Sign the token with the secret
	//secret := "a secret message"
	tokenString, err := token.SignedString([]byte(app.config.jwt.secret))
	if err != nil {
		return
	}

	// Send JWT token back to the client
	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": tokenString}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// validateUser is a helper function to validate user. Should go to DB and validate user, hardcoded for now
func (app *application) validateUser(username string, password string) bool {
	return username != "john@companyservice.io" || password != "doe"
}
