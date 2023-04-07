package main

import (
	"database/sql"
	"github.com/google/uuid"
	"mborgnolo/companyservice/internal/data"
	"mborgnolo/companyservice/internal/validator"
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

func (app *application) CreateCompanyHandler(writer http.ResponseWriter, request *http.Request) {
	var uuid uuid.UUID
	var input struct {
		Name        string         `json:"name"`
		Description sql.NullString `json:"description"`
		Employees   int            `json:"employees"`
		Registered  bool           `json:"registered"`
		Type        string         `json:"type"`
	}
	err := app.readJSON(request, &input)
	if err != nil {
		app.badRequestResponse(writer, request, err)
		return
	}
	company := &data.Company{
		Name:        input.Name,
		Description: input.Description,
		Employees:   input.Employees,
		Registered:  input.Registered,
		Type:        input.Type,
	}

	v := validator.New()

	if data.ValidateCompany(v, company); !v.IsValid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}
	uuid, err = app.company.CreateCompany(company)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}
	err = app.writeJSON(writer, http.StatusCreated, envelope{"id": uuid}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
