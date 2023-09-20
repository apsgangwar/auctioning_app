package svc

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"demand/objects"

	"github.com/google/uuid"
	"github.com/DATA-DOG/go-sqlmock"
)

type MockDB struct {
	sql.DB
}

func (mdb *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	s := args[0].(string)
	if s == "valid-id" {
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow("valid-id", "valid-name")
		mockRow := new(sqlmock.Rows)
		mockRow.RowError(1, rows.Err())
		return mockRow.Row()
	}
	return &sql.Row{}
}

func TestGetOne_8029431e34(t *testing.T) {
	DB = &MockDB{}

	req, err := http.NewRequest("GET", "/bidders/valid-id", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bidder, err := GetOne(w, r, "valid-id")
		if err != nil {
			t.Error("Expected no error for valid ID, got", err)
			return
		}
		if bidder.Id != "valid-id" || bidder.Name != "valid-name" {
			t.Error("Unexpected bidder data for valid ID")
			return
		}
		t.Log("Test case 1: Passed")
	})
	handler.ServeHTTP(rr, req)

	req, err = http.NewRequest("GET", "/bidders/invalid-id", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := GetOne(w, r, "invalid-id")
		if err == nil {
			t.Error("Expected error for invalid ID, got nil")
			return
		}
		t.Log("Test case 2: Passed")
	})
	handler.ServeHTTP(rr, req)
}
