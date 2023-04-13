package data

import (
	"database/sql"
	"github.com/google/uuid"
	"mborgnolo/companyservice/internal/validator"
)

// CompanyDescription is a custom type that wraps a string and implements the
// sql.Scanner interface. This allows us to store a NULL
// value in the database if the user doesn't provide a description for the
// company. It also implements the json.Marshaler and json.Unmarshaler interfaces
// so that we can control how the value is encoded and decoded when it is
// sent to and received from the client.
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

// Company represents a company in the companyservice application.
type Company struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Description CompanyDescription `json:"description"`
	Employees   int                `json:"employees"`
	//FIXME: This should be an optional field
	Registered *bool  `json:"registered"`
	Type       string `json:"type"`
}

// ValidateCompany runs validation checks on the company data.
func ValidateCompany(v *validator.Validator, company *Company) {
	v.Check(company.Name != "", "name", "is required")
	v.Check(len(company.Description.String) < 3000, "description", "must be less than 3000 characters")
	v.Check(company.Employees > 0, "employees", "must be greater than zero")
	//FIXME: This should be an optional field
	v.Check(company.Registered != nil, "registered", "is required")
	v.Check(company.Type != "", "type", "is required")
	v.Check(validateCompanyType(company.Type), "type", "must be one of: Corporations, NonProfit, Cooperative , Sole Proprietorship")
}

func validateCompanyType(t string) bool {
	validTypes := []string{"Corporations", "NonProfit", "Cooperative", "Sole Proprietorship"}
	for _, v := range validTypes {
		if v == t {
			return true
		}
	}
	return false
}

// CompanyModel wraps the sql.DB connection pool.
type CompanyModel struct {
	DB *sql.DB
}

// GetCompany returns a single company based on the ID provided.
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

// CreateCompany inserts a new company record in the database.
func (m *CompanyModel) CreateCompany(company *Company) (uuid.UUID, error) {
	newUUID := uuid.New()
	query := `INSERT INTO company ("id", "name", "description", "employees", "registered", "type") VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := m.DB.Exec(query, newUUID, company.Name, company.Description.String, company.Employees, company.Registered, company.Type)
	if err != nil {
		return uuid.Nil, err
	}
	return newUUID, nil
}

// DeleteCompany deletes a company record from the database.
func (m *CompanyModel) DeleteCompany(id uuid.UUID) error {
	query := `DELETE FROM company WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateCompany updates a company record in the database.
func (m *CompanyModel) UpdateCompany(company *Company) error {
	query := `UPDATE company SET name = COALESCE($1,name), description = COALESCE($2,description), employees = COALESCE($3,employees), registered = COALESCE($4,registered), type = COALESCE($5,type) WHERE id = $6`
	result, err := m.DB.Exec(query, company.Name, company.Description.String, company.Employees, company.Registered, company.Type, company.ID)
	if err != nil {
		return err
	}
	i, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if i == 0 {
		return ErrRecordNotFound
	}

	return nil
}
