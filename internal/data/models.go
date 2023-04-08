package data

import (
	"database/sql"
	"encoding/json"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type NullString sql.NullString

func (x *NullString) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(x.String)
}

// NewCompanyModel returns a new CompanyModel.
func NewCompanyModel(db *sql.DB) *CompanyModel {
	return &CompanyModel{DB: db}
}
