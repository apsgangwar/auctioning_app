package svc

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"demand/objects"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func TestGetAll_4ef88fc27e(t *testing.T) {
	// TODO: Setup Test DB Connection
	DB, _ = sql.Open("mysql", "user:password@/dbname")

	// Test case 1: Successful data retrieval
	req, err := http.NewRequest("GET", "/getAll", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAll)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	var response objects.GetAllBidderResp
	json.Unmarshal(rr.Body.Bytes(), &response)
	if len(response.Data) == 0 {
		t.Errorf("handler returned unexpected body: got empty result")
	} else {
		t.Log("TestGetAll_4ef88fc27e: Test case 1 passed")
	}

	// Test case 2: Failure due to DB connection error
	// TODO: Close DB connection before running this test case
	DB.Close()
	req, err = http.NewRequest("GET", "/getAll", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetAll)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	} else {
		t.Log("TestGetAll_4ef88fc27e: Test case 2 passed")
	}
}
