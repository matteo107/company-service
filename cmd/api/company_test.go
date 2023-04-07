package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCompany(t *testing.T) {
	app := newTestApplication(t)
	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/v1/company/dc152cf7-cc4b-4555-8d4c-1878e5b9262c", http.StatusOK, []byte("company")},
		{"Non-existent ID", "/v1/company/5f001b5d-8cd1-4f90-8a6a-5164adee43b5", http.StatusNotFound, nil},
		{"Empty ID", "/v1/company/", http.StatusNotFound, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs, err := ts.Client().Get(ts.URL + tt.urlPath)
			if err != nil {
				t.Fatal(err)
			}
			defer rs.Body.Close()
			body, err := io.ReadAll(rs.Body)
			if err != nil {
				t.Fatal(err)
			}

			if rs.StatusCode != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, rs.StatusCode)
			}
			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q; got %q", tt.wantBody, body)
			}
		})
	}
}

func TestCreateCompany(t *testing.T) {
	app := newTestApplication(t)
	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	tests := []struct {
		name        string
		urlPath     string
		bodyRequest io.Reader
		wantCode    int
		wantBody    []byte
	}{
		{"Valid ID",
			"/v1/company/",
			bytes.NewReader([]byte(`{"name":"AWS","employees":1000,"registered":true,"type":"Corporate"}`)),
			http.StatusCreated,
			[]byte("id")},
		{"Mandatory data missing",
			"/v1/company/",
			bytes.NewReader([]byte(`{"name":"AWS","registered":true,"type":"Corporate"}`)),
			http.StatusUnprocessableEntity,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs, err := ts.Client().Post(ts.URL+tt.urlPath, "application/json", tt.bodyRequest)
			if err != nil {
				t.Fatal(err)
			}
			defer rs.Body.Close()
			body, err := io.ReadAll(rs.Body)
			if err != nil {
				t.Fatal(err)
			}

			if rs.StatusCode != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, rs.StatusCode)
			}
			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q; got %q", tt.wantBody, body)
			}
		})
	}
}

func TestDeleteCompany(t *testing.T) {
	app := newTestApplication(t)
	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/v1/company/dc152cf7-cc4b-4555-8d4c-1878e5b9262c", http.StatusOK, nil},
		{"Non-existent ID", "/v1/company/5f001b5d-8cd1-4f90-8a6a-5164adee43b5", http.StatusNotFound, nil},
		{"Empty ID", "/v1/company/", http.StatusNotFound, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, ts.URL+tt.urlPath, nil)
			if err != nil {
				t.Fatal(err)
			}
			rs, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer rs.Body.Close()
			body, err := io.ReadAll(rs.Body)
			if err != nil {
				t.Fatal(err)
			}

			if rs.StatusCode != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, rs.StatusCode)
			}
			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q; got %q", tt.wantBody, body)
			}
		})
	}
}
