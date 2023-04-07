package mocks

import (
	"github.com/google/uuid"
	"mborgnolo/companyservice/internal/data"
)

var mockCompany = &data.Company{
	ID:          uuid.MustParse("dc152cf7-cc4b-4555-8d4c-1878e5b9262c"),
	Name:        "Test Company",
	Description: data.CompanyDescription{"Test Company Description", true},
	Registered:  true,
	Type:        "Corporate",
}

type CompanyModel struct{}

func (t *CompanyModel) GetCompany(id uuid.UUID) (*data.Company, error) {
	if id.String() == mockCompany.ID.String() {
		return mockCompany, nil
	}
	return nil, data.ErrRecordNotFound
}

func (t *CompanyModel) CreateCompany(company *data.Company) (uuid.UUID, error) {
	return uuid.MustParse("dc152cf7-cc4b-4555-8d4c-1878e5b9262c"), nil
}

func (t *CompanyModel) DeleteCompany(id uuid.UUID) error {
	return nil
}

func (t *CompanyModel) UpdateCompany(company *data.Company) error {
	return nil
}
