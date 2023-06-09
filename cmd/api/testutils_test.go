package main

import (
	"log"
	"mborgnolo/companyservice/internal/data"
	"mborgnolo/companyservice/internal/mocks"
	"os"
	"testing"
)

// newTestApplication returns an instance of application configured for testing
func newTestApplication(t *testing.T) *application {

	return &application{
		config:    config{env: "test"},
		logger:    log.New(os.Stdout, "", log.Ldate|log.Ltime),
		company:   &mocks.CompanyModel{},
		eventChan: make(chan data.EventRecord, 100),
	}
}
