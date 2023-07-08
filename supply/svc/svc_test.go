package svc_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"supply/objects"
	"supply/svc"
)

func TestGetAll(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Set up the mock expectations
	rows := sqlmock.NewRows([]string{"id", "name", "base_price"}).
		AddRow("1", "Space 1", 10.0).
		AddRow("2", "Space 2", 15.0)
	mock.ExpectQuery("SELECT (.+) FROM ad_space_auction.space").
		WillReturnRows(rows)

	// Set the mock DB in the svc package
	svc.DB = db

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/space", nil)
	assert.NoError(t, err)

	// Create a test HTTP response recorder
	res := httptest.NewRecorder()

	// Call the GetAll function
	resp, err := svc.GetAll(res, req)
	assert.NoError(t, err)

	// Assert the expected response
	expected := objects.GetAllSpaceResp{
		Data: []objects.SpaceEntity{
			{Id: "1", Name: "Space 1", BasePrice: 10.0},
			{Id: "2", Name: "Space 2", BasePrice: 15.0},
		},
	}
	assert.Equal(t, expected, resp)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOne(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Set up the mock expectations
	rows := sqlmock.NewRows([]string{"id", "name", "base_price"}).
		AddRow("1", "Space 1", 10.0)
	mock.ExpectQuery("SELECT (.+) FROM ad_space_auction.space WHERE id = ?").
		WithArgs("1").
		WillReturnRows(rows)

	// Set the mock DB in the svc package
	svc.DB = db

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/space/1", nil)
	assert.NoError(t, err)

	// Create a test HTTP response recorder
	res := httptest.NewRecorder()

	// Call the GetAll function
	resp, err := svc.GetOne(res, req, "1")
	assert.NoError(t, err)

	// Assert the expected response
	expected := objects.SpaceEntity{Id: "1", Name: "Space 1", BasePrice: 10.0}
	assert.Equal(t, expected, resp)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Set up the mock expectations
	mock.ExpectExec("INSERT INTO ad_space_auction.space").
		WithArgs(sqlmock.AnyArg(), "Test Space", 10.0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Set the mock DB in the svc package
	svc.DB = db

	// Create a test space entity
	space := objects.SpaceEntity{
		Name:      "Test Space",
		BasePrice: 10.0,
	}

	// Convert the space entity to JSON
	body, err := json.Marshal(space)
	assert.NoError(t, err)

	// Create a test HTTP request with the JSON body
	req, err := http.NewRequest("POST", "/space", bytes.NewBuffer(body))
	assert.NoError(t, err)

	// Create a test HTTP response recorder
	res := httptest.NewRecorder()

	// Call the Create function
	resp, err := svc.Create(res, req)
	assert.NoError(t, err)

	// Assert the expected response
	assert.NotEmpty(t, resp.Id)
	assert.Equal(t, "Test Space", resp.Name)
	assert.Equal(t, 10.0, resp.BasePrice)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Set up the mock expectations
	rows := sqlmock.NewRows([]string{"id", "name", "base_price"}).
		AddRow("1", "Space 1", 10.0)
	mock.ExpectQuery("SELECT (.+) FROM ad_space_auction.space WHERE id = ?").
		WithArgs("1").
		WillReturnRows(rows)
	mock.ExpectExec("UPDATE ad_space_auction.space SET").
		WithArgs("Updated Space", 20.0, "1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Set the mock DB in the svc package
	svc.DB = db

	// Create a test space entity
	space := objects.SpaceEntity{
		Id:        "1",
		Name:      "Updated Space",
		BasePrice: 20.0,
	}

	// Convert the space entity to JSON
	body, err := json.Marshal(space)
	assert.NoError(t, err)

	// Create a test HTTP request with the JSON body
	req, err := http.NewRequest("PATCH", "/space/1", bytes.NewBuffer(body))
	assert.NoError(t, err)

	// Create a test HTTP response recorder
	res := httptest.NewRecorder()

	// Call the Update function
	resp, err := svc.Update(res, req, "1")
	assert.NoError(t, err)

	// Assert the expected response
	assert.Equal(t, "Updated Space", resp.Name)
	assert.Equal(t, 20.0, resp.BasePrice)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Set up the mock expectations
	rows := sqlmock.NewRows([]string{"id", "name", "base_price"}).
		AddRow("1", "Space 1", 10.0)
	mock.ExpectQuery("SELECT (.+) FROM ad_space_auction.space WHERE id = ?").
		WithArgs("123").
		WillReturnError(errors.New("not found"))
	mock.ExpectQuery("SELECT (.+) FROM ad_space_auction.space WHERE id = ?").
		WithArgs("1").
		WillReturnRows(rows)
	mock.ExpectExec("DELETE FROM ad_space_auction.space").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Set the mock DB in the svc package
	svc.DB = db

	// Create a test HTTP request
	req, err := http.NewRequest("DELETE", "/space/123", nil)
	assert.NoError(t, err)

	// Create a test HTTP response recorder
	res := httptest.NewRecorder()

	// Call the Delete function
	resp, err := svc.Delete(res, req, "123")
	assert.Error(t, err)
	assert.Equal(t, objects.SpaceEntity{}, resp)

	// Call the Delete function again with different id
	resp, err = svc.Delete(res, req, "1")
	assert.NoError(t, err)

	// Assert the expected response
	assert.Equal(t, "1", resp.Id)
	assert.Equal(t, "Space 1", resp.Name)
	assert.Equal(t, 10.0, resp.BasePrice)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
