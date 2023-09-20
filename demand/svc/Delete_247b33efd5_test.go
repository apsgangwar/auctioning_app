package svc

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"demand/objects"
)

type MockDB struct {
	ExecFunc func(query string, args ...interface{}) (sql.Result, error)
}

var (
	// GetOne is a function that mocks the GetOne function
	GetOne func(w http.ResponseWriter, r *http.Request, id string) (objects.BidderEntity, error)

	// DB is a mock database
	DB *MockDB
)

func TestDelete_247b33efd5(t *testing.T) {
	GetOne = func(w http.ResponseWriter, r *http.Request, id string) (objects.BidderEntity, error) {
		return objects.BidderEntity{Id: "test-id"}, nil
	}

	DB = &MockDB{
		ExecFunc: func(query string, args ...interface{}) (sql.Result, error) {
			return nil, nil
		},
	}

	req, err := http.NewRequest("DELETE", "/delete/test-id", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Delete)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Delete handler responded with the wrong status code")

	GetOne = func(w http.ResponseWriter, r *http.Request, id string) (objects.BidderEntity, error) {
		return objects.BidderEntity{}, errors.New("get one error")
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Delete handler responded with the wrong status code when GetOne returns an error")

	DB = &MockDB{
		ExecFunc: func(query string, args ...interface{}) (sql.Result, error) {
			return nil, errors.New("db exec error")
		},
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Delete handler responded with the wrong status code when DB.Exec returns an error")
}
