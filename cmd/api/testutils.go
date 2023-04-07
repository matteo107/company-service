package main

import (
	"log"
	"mborgnolo/companyservice/internal/mocks"
	"os"
	"testing"
)

func newTestApplication(t *testing.T) *application {

	return &application{
		config:  config{},
		logger:  log.New(os.Stdout, "", log.Ldate|log.Ltime),
		company: &mocks.CompanyModel{},
	}
}
