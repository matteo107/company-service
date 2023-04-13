package data

import (
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"reflect"
	"testing"
)

func TestCompanyModelGet(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}
	tests := []struct {
		name        string
		companyID   uuid.UUID
		wantCompany *Company
		wantError   error
	}{
		{
			name:      "Get Company",
			companyID: uuid.MustParse("f1203d76-0491-47fe-9640-0aeda76ad3f6"),
			wantCompany: &Company{
				ID:          uuid.MustParse("f1203d76-0491-47fe-9640-0aeda76ad3f6"),
				Name:        "Company One",
				Description: CompanyDescription{String: "Description for company one", Valid: true},
				Employees:   100,
				Registered:  boolPtr(true),
				Type:        "Corporations",
			},
			wantError: nil,
		},
		{
			name:        "Non-existent ID",
			companyID:   uuid.MustParse("e2d3253c-3e65-4516-9318-d013fde56dca"),
			wantCompany: nil,
			wantError:   ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, teardown := newTestDB(t)
			defer teardown()

			c := CompanyModel{db}

			company, err := c.GetCompany(tt.companyID)
			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}
			if !reflect.DeepEqual(company, tt.wantCompany) {
				t.Errorf("want %v; got %v", tt.wantCompany, company)
			}
		})
	}

}

func TestCompanyModelUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}
	tests := []struct {
		name          string
		companyID     uuid.UUID
		updateCompany *Company
		wantCompany   *Company
		wantError     error
	}{
		{
			name:      "Update Company",
			companyID: uuid.MustParse("f1203d76-0491-47fe-9640-0aeda76ad3f6"),
			updateCompany: &Company{
				ID:          uuid.MustParse("f1203d76-0491-47fe-9640-0aeda76ad3f6"),
				Name:        "Company One",
				Description: CompanyDescription{String: "Description for company one", Valid: true},
				Employees:   2,
				Registered:  boolPtr(true),
				Type:        "Corporations",
			},
			wantCompany: &Company{
				ID:          uuid.MustParse("f1203d76-0491-47fe-9640-0aeda76ad3f6"),
				Name:        "Company One",
				Description: CompanyDescription{String: "Description for company one", Valid: true},
				Employees:   2,
				Registered:  boolPtr(true),
				Type:        "Corporations",
			},
			wantError: nil,
		},
		{
			name:      "Non-existent ID",
			companyID: uuid.MustParse("e2d3253c-3e65-4516-9318-d013fde56dca"),
			updateCompany: &Company{
				ID:          uuid.MustParse("e2d3253c-3e65-4516-9318-d013fde56dca"),
				Description: CompanyDescription{String: "Description for company one", Valid: true},
			},
			wantCompany: nil,
			wantError:   ErrRecordNotFound,
		},
		{
			name:      "Update Company with empty Registered field",
			companyID: uuid.MustParse("f1203d76-0491-47fe-9640-0aeda76ad3f6"),
			updateCompany: &Company{
				ID:          uuid.MustParse("f1203d76-0491-47fe-9640-0aeda76ad3f6"),
				Name:        "Company One",
				Description: CompanyDescription{String: "Description for company one", Valid: true},
				Employees:   2,
				Type:        "Corporations",
			},
			wantCompany: &Company{
				ID:          uuid.MustParse("f1203d76-0491-47fe-9640-0aeda76ad3f6"),
				Name:        "Company One",
				Description: CompanyDescription{String: "Description for company one", Valid: true},
				Employees:   2,
				Registered:  boolPtr(true),
				Type:        "Corporations",
			},
			wantError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()
			c := CompanyModel{db}
			err := c.UpdateCompany(tt.updateCompany)
			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}
			companyFromDb, err := c.GetCompany(tt.updateCompany.ID)
			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}
			if !reflect.DeepEqual(companyFromDb, tt.wantCompany) {
				t.Errorf("want %v; got %v", tt.wantCompany, companyFromDb)
			}
		})
	}

}

func TestCompanyModelCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres: skipping integration test")
	}
	tests := []struct {
		name      string
		company   *Company
		wantError error
	}{
		{
			name: "Create Company",
			company: &Company{
				Name:        "Company Two",
				Description: CompanyDescription{String: "Description for company one", Valid: true},
				Employees:   500,
				Registered:  boolPtr(true),
				Type:        "NonProfit",
			},
			wantError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, teardown := newTestDB(t)
			defer teardown()

			c := CompanyModel{db}
			UUID, err := c.CreateCompany(tt.company)

			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}
			companyFromDb, err := c.GetCompany(UUID)
			// Set the UUID of the company to the one returned from the database, since it's generated
			// by the database
			tt.company.ID = UUID
			if !reflect.DeepEqual(companyFromDb, tt.company) {
				t.Errorf("want %v; got %v", tt.company, companyFromDb)
			}
		})
	}
}
