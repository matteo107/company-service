package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

// This function is used to create a new router instance and register all the application routes.
// It also registers a middleware function (app.authenticate) that will be called before any of the
// handlers used for mutating operation are executed.
func (app *application) routes() http.Handler {
	router := httprouter.New()
	standardMiddleware := alice.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/company/:id", app.GetCompanyHandler)
	router.Handler(http.MethodPost, "/v1/company", standardMiddleware.Append(app.authenticate).ThenFunc(app.CreateCompanyHandler))
	router.Handler(http.MethodPatch, "/v1/company/:id", standardMiddleware.Append(app.authenticate).ThenFunc(app.UpdateCompanyHandler))
	router.Handler(http.MethodDelete, "/v1/company/:id", standardMiddleware.Append(app.authenticate).ThenFunc(app.DeleteCompanyHandler))
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	return standardMiddleware.Then(router)
}
