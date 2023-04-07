package main

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) readJSON(r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) readIDParam(r *http.Request) (uuid.UUID, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return uuid.Nil, errors.New("invalid id parameter")
	}
	return id, nil

}
