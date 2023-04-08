package main

import (
	"context"
	"errors"
	"github.com/pascaldekloe/jwt"
	"mborgnolo/companyservice/internal/data"
	"net/http"
	"strings"
	"time"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if !claims.Valid(time.Now()) {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if claims.Issuer != "api.companyservice.io" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		if !claims.AcceptAudience("api.companyservice.io") {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		// hardcoded for now
		if claims.Subject != ClaimSubject {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, claims.Subject)
		next.ServeHTTP(w, r)
	})
}

func (app *application) contextSetUser(r *http.Request, data string) *http.Request {
	user := userContextKey("user")
	ctx := context.WithValue(r.Context(), user, data)
	return r.WithContext(ctx)
}

type userContextKey string
