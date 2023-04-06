package main

import (
	"mborgnolo/companyservice/internal/data"
	"net/http"
)

func (app *application) GetCompanyHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := app.readIDParam(request)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}
	company, err := app.company.GetCompany(id)
	if err != nil {
		switch err {
		case data.ErrRecordNotFound:
			app.notFoundResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}
	err = app.writeJSON(writer, http.StatusOK, envelope{"company": company}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
