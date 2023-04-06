package main

import "net/http"
import "github.com/julienschmidt/httprouter"

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	//router.HandlerFunc(http.MethodGet, "/v1/company/:id", app.GetCompanyHandler)
	//router.HandlerFunc(http.MethodPost, "/v1/company", app.CreateCompanyHandler)
	//router.HandlerFunc(http.MethodPatch, "/v1/company/:id", app.UpdateCompanyHandler)
	//router.HandlerFunc(http.MethodDelete, "/v1/company/:id", app.DeleteCompanyHandler)
	return router
}
