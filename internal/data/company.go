package data

import (
	"database/sql"
	"github.com/google/uuid"
)

type Company struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Employees   int       `json:"employees"`
	Registered  bool      `json:"registered"`
	Type        string    `json:"type"`
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

func (m *CompanyModel) InsertCompany(company *Company) error {
	return nil
}

func (m *CompanyModel) DeleteCompany(id uuid.UUID) error {
	return nil
}

func (m *CompanyModel) UpdateCompany(company *Company) error {
	return nil
}
