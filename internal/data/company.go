package data

import (
	"database/sql"
	"github.com/google/uuid"
	"mborgnolo/companyservice/internal/validator"
)

type CompanyDescription sql.NullString

func (cd *CompanyDescription) MarshalJSON() ([]byte, error) {
	if !cd.Valid {
		return []byte("null"), nil
	}
	return []byte(`"` + cd.String + `"`), nil
}

func (cd *CompanyDescription) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		cd.Valid = false
		return nil
	}
	cd.Valid = true
	tmp := string(data)
	cd.String = tmp[1 : len(tmp)-1]
	return nil
}

func (cd *CompanyDescription) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		cd.String = v
		cd.Valid = true
	case []byte:
		cd.String = string(v)
		cd.Valid = true
	case nil:
		cd.Valid = false
	}
	return nil
}

type Company struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Description CompanyDescription `json:"description"`
	Employees   int                `json:"employees"`
	Registered  bool               `json:"registered"`
	Type        string             `json:"type"`
}

func ValidateCompany(v *validator.Validator, company *Company) {
	v.Check(company.Name != "", "name", "is required")
	v.Check(len(company.Description.String) < 3000, "description", "must be less than 3000 characters")
	v.Check(company.Employees > 0, "employees", "must be greater than zero")
	v.Check(company.Registered, "registered", "is required")
	v.Check(company.Type != "", "type", "is required")
}

type CompanyModel struct {
	DB *sql.DB
}

func (m *CompanyModel) GetCompany(id uuid.UUID) (*Company, error) {
	query := `SELECT "id", "name", "description", "employees", "registered", "type" FROM company WHERE id = $1`
	row := m.DB.QueryRow(query, id)
	company := &Company{}
	err := row.Scan(&company.ID, &company.Name, &company.Description, &company.Employees, &company.Registered, &company.Type)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return company, nil
}

func (m *CompanyModel) CreateCompany(company *Company) (uuid.UUID, error) {
	newUUID := uuid.New()
	query := `INSERT INTO company ("id", "name", "description", "employees", "registered", "type") VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := m.DB.Exec(query, newUUID, company.Name, company.Description.String, company.Employees, company.Registered, company.Type)
	if err != nil {
		return uuid.Nil, err
	}
	return newUUID, nil
}

func (m *CompanyModel) DeleteCompany(id uuid.UUID) error {
	query := `DELETE FROM company WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *CompanyModel) UpdateCompany(company *Company) error {
	query := `UPDATE company SET "name" = $1, "description" = $2, "employees" = $3, "registered" = $4, "type" = $5 WHERE id = $6`
	_, err := m.DB.Exec(query, company.Name, company.Description.String, company.Employees, company.Registered, company.Type, company.ID)
	if err != nil {
		return err
	}
	return nil
}
