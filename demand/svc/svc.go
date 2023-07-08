package svc

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"demand/objects"

	"github.com/google/uuid"
)

var DB *sql.DB

func GetAll(w http.ResponseWriter, r *http.Request) (objects.GetAllBidderResp, error) {
	var resp = objects.GetAllBidderResp{}

	rows, err := DB.Query("SELECT * FROM ad_space_auction.bidder")
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	var bidders = []objects.BidderEntity{}
	for rows.Next() {
		var bidder objects.BidderEntity
		err := rows.Scan(&bidder.Id, &bidder.Name)
		if err != nil {
			return resp, err
		}
		bidders = append(bidders, bidder)
	}

	if err := rows.Err(); err != nil {
		return resp, err
	}

	resp.Data = bidders
	return resp, err
}

func GetOne(w http.ResponseWriter, r *http.Request, id string) (objects.BidderEntity, error) {
	row := DB.QueryRow("SELECT * FROM ad_space_auction.bidder WHERE id = ?", id)

	var bidder objects.BidderEntity
	err := row.Scan(&bidder.Id, &bidder.Name)
	if err == sql.ErrNoRows {
		http.Error(w, "entity not found", http.StatusNotFound)
		return bidder, err
	}
	if err != nil {
		return bidder, err
	}

	return bidder, err
}

func Create(w http.ResponseWriter, r *http.Request) (objects.BidderEntity, error) {
	var bidder objects.BidderEntity

	err := json.NewDecoder(r.Body).Decode(&bidder)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return bidder, err
	}

	if bidder.Name == "" {
		err = errors.New("name shouldn't be empty")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return bidder, err
	}

	id := uuid.NewString()
	bidder.Id = id

	_, err = DB.Exec("INSERT INTO ad_space_auction.bidder (id, name) VALUES (?, ?)", bidder.Id, bidder.Name)
	if err != nil {
		return bidder, err
	}

	return bidder, err
}

func Update(w http.ResponseWriter, r *http.Request, id string) (objects.BidderEntity, error) {
	var bidder objects.BidderEntity

	err := json.NewDecoder(r.Body).Decode(&bidder)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return bidder, err
	}

	if bidder.Name == "" {
		err = errors.New("name shouldn't be empty")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return bidder, err
	}

	bidder.Id = id

	_, err = GetOne(w, r, id)
	if err != nil {
		return bidder, err
	}

	_, err = DB.Exec("UPDATE ad_space_auction.bidder SET name = ? WHERE id = ?", bidder.Name, bidder.Id)
	if err != nil {
		return bidder, err
	}

	return bidder, err
}

func Delete(w http.ResponseWriter, r *http.Request, id string) (objects.BidderEntity, error) {
	bidder, err := GetOne(w, r, id)
	if err != nil {
		return bidder, err
	}

	_, err = DB.Exec("DELETE FROM ad_space_auction.bidder WHERE id = ?", bidder.Id)
	if err != nil {
		return bidder, err
	}

	return bidder, err
}
