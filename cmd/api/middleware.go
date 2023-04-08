package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

// authenticate is a middleware function which will be used to authenticate requests
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if testing, skip authentication
		if app.config.env == "test" {
			next.ServeHTTP(w, r)
			return
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Verify JWT token
		tokenString := authHeader[len("Bearer "):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(app.config.jwt.secret), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				app.invalidAuthenticationTokenResponse(w, r)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !token.Valid {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// If JWT token is valid, call next handler
		next.ServeHTTP(w, r)

	})
}
