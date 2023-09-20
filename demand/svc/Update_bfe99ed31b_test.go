package svc

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"demand/objects"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func TestUpdate_bfe99ed31b(t *testing.T) {
	// Test case 1: Successful update
	{
		bidder := objects.BidderEntity{
			Name: "TestBidder",
		}
		id := uuid.New().String()
		bidderJSON, _ := json.Marshal(bidder)
		req, err := http.NewRequest("POST", "/update", bytes.NewBuffer(bidderJSON))
		if err != nil {
			t.Fatal(err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Update)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var updatedBidder objects.BidderEntity
		err = json.NewDecoder(rr.Body).Decode(&updatedBidder)
		if err != nil {
			t.Fatal(err)
		}

		if updatedBidder.Name != bidder.Name {
			t.Errorf("handler returned unexpected body: got %v want %v", updatedBidder.Name, bidder.Name)
		}
	}

	// Test case 2: Empty bidder name
	{
		bidder := objects.BidderEntity{
			Name: "",
		}
		id := uuid.New().String()
		bidderJSON, _ := json.Marshal(bidder)
		req, err := http.NewRequest("POST", "/update", bytes.NewBuffer(bidderJSON))
		if err != nil {
			t.Fatal(err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Update)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expectedError := "name shouldn't be empty"
		if rr.Body.String() != expectedError {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedError)
		}
	}
}
