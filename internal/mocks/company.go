package mocks

import (
	"github.com/google/uuid"
	"mborgnolo/companyservice/internal/data"
)

func boolPtr(b bool) *bool {
	return &b
}

// MockCompany is a mock company used for testing.
var mockCompany = &data.Company{
	ID:          uuid.MustParse("dc152cf7-cc4b-4555-8d4c-1878e5b9262c"),
	Name:        "Test Company",
	Description: data.CompanyDescription{String: "Test Company Description", Valid: true},
	Registered:  boolPtr(true),
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
	if id.String() == mockCompany.ID.String() {
		return nil
	}
	return data.ErrRecordNotFound
}

func (t *CompanyModel) UpdateCompany(company *data.Company) error {
	if company.ID.String() == mockCompany.ID.String() {
		return nil
	}
	return data.ErrRecordNotFound
}
