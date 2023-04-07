package main

import (
	"github.com/google/uuid"
	"mborgnolo/companyservice/internal/data"
	"mborgnolo/companyservice/internal/validator"
	"net/http"
	"time"
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
	var UUID uuid.UUID
	var input struct {
		Name        string                  `json:"name"`
		Description data.CompanyDescription `json:"description"`
		Employees   int                     `json:"employees"`
		Registered  bool                    `json:"registered"`
		Type        string                  `json:"type"`
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
	UUID, err = app.company.CreateCompany(company)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}
	err = app.writeJSON(writer, http.StatusCreated, envelope{"id": UUID}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
	app.eventChan <- data.EventRecord{
		ID:        UUID,
		Type:      data.CompanyCreated,
		TimeStamp: time.Now().UTC(),
	}
}

func (app *application) DeleteCompanyHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := app.readIDParam(request)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}
	err = app.company.DeleteCompany(id)
	if err != nil {
		switch err {
		case data.ErrRecordNotFound:
			app.notFoundResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}
	err = app.writeJSON(writer, http.StatusOK, envelope{"id": id}, nil)
	app.eventChan <- data.EventRecord{
		ID:        id,
		Type:      data.CompanyDeleted,
		TimeStamp: time.Now().UTC(),
	}
}

func (app *application) UpdateCompanyHandler(writer http.ResponseWriter, request *http.Request) {
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
	var input struct {
		Name        *string                 `json:"name"`
		Description data.CompanyDescription `json:"description"`
		Employees   *int                    `json:"employees"`
		Registered  *bool                   `json:"registered"`
		Type        *string                 `json:"type"`
	}
	err = app.readJSON(request, &input)
	if err != nil {
		app.badRequestResponse(writer, request, err)
		return
	}
	if input.Name != nil {
		company.Name = *input.Name
	}
	if input.Description.Valid {
		company.Description.String = input.Description.String
		company.Description.Valid = input.Description.Valid
	}
	if input.Employees != nil {
		company.Employees = *input.Employees
	}
	if input.Registered != nil {
		company.Registered = *input.Registered
	}
	if input.Type != nil {
		company.Type = *input.Type
	}
	v := validator.New()
	if data.ValidateCompany(v, company); !v.IsValid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}
	err = app.company.UpdateCompany(company)

	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}
	err = app.writeJSON(writer, http.StatusOK, envelope{"company": company}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
	app.eventChan <- data.EventRecord{
		ID:        id,
		Type:      data.CompanyUpdated,
		TimeStamp: time.Now().UTC(),
	}
}
