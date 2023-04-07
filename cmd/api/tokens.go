package main

import (
	"github.com/pascaldekloe/jwt"
	"net/http"
	"time"
)

// hardcoded for now
const ClaimSubject = "user@companyservice.io"

func (app *application) createAuthenticationTokenHandler(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(request, &input)
	if err != nil {
		app.badRequestResponse(writer, request, err)
		return
	}

	var claims jwt.Claims

	claims.Subject = ClaimSubject
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "api.companyservice.io"
	claims.Audiences = []string{"api.companyservice.io"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}
	err = app.writeJSON(writer, http.StatusCreated, envelope{"authentication_token": string(jwtBytes)}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
