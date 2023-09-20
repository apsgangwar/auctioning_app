package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInitialize_9db7c49ca0(t *testing.T) {
	Initialize()

	// Test case 1: Check if "/bidder" route is working
	req, err := http.NewRequest("GET", "/bidder", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleService)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Test case 2: Check if "/bidder/" route is working
	req, err = http.NewRequest("GET", "/bidder/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(HandleService)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
