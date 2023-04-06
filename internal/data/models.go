package data

import (
	"database/sql"
	"encoding/json"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Trades CompanyModel
}

type NullString sql.NullString

func (x *NullString) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(x.String)
}

func NewCompanyModel(db *sql.DB) *CompanyModel {
	return &CompanyModel{DB: db}
}
