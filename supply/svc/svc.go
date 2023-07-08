package svc

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"supply/objects"

	"github.com/google/uuid"
)

var DB *sql.DB

func GetAll(w http.ResponseWriter, r *http.Request) (objects.GetAllSpaceResp, error) {
	var resp = objects.GetAllSpaceResp{}

	rows, err := DB.Query("SELECT * FROM ad_space_auction.space")
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	var spaces = []objects.SpaceEntity{}
	for rows.Next() {
		var space objects.SpaceEntity
		err := rows.Scan(&space.Id, &space.Name, &space.BasePrice)
		if err != nil {
			return resp, err
		}
		spaces = append(spaces, space)
	}

	if err := rows.Err(); err != nil {
		return resp, err
	}

	resp.Data = spaces
	return resp, err
}

func GetOne(w http.ResponseWriter, r *http.Request, id string) (objects.SpaceEntity, error) {
	row := DB.QueryRow("SELECT * FROM ad_space_auction.space WHERE id = ?", id)

	var space objects.SpaceEntity
	err := row.Scan(&space.Id, &space.Name, &space.BasePrice)
	if err == sql.ErrNoRows {
		http.Error(w, "entity not found", http.StatusNotFound)
		return space, err
	}
	if err != nil {
		return space, err
	}

	return space, err
}

func Create(w http.ResponseWriter, r *http.Request) (objects.SpaceEntity, error) {
	var space objects.SpaceEntity

	err := json.NewDecoder(r.Body).Decode(&space)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return space, err
	}

	if space.Name == "" {
		err = errors.New("name shouldn't be empty")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return space, err
	} else if space.BasePrice <= 0 {
		err = errors.New("base_price should be a positive number")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return space, err
	}

	id := uuid.NewString()
	space.Id = id

	_, err = DB.Exec("INSERT INTO ad_space_auction.space (id, name, base_price) VALUES (?, ?, ?)", space.Id, space.Name, space.BasePrice)
	if err != nil {
		return space, err
	}

	return space, err
}

func Update(w http.ResponseWriter, r *http.Request, id string) (objects.SpaceEntity, error) {
	var space objects.SpaceEntity

	err := json.NewDecoder(r.Body).Decode(&space)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return space, err
	}

	if space.Name == "" {
		err = errors.New("name shouldn't be empty")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return space, err
	} else if space.BasePrice <= 0 {
		err = errors.New("base_price should be a positive number")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return space, err
	}

	space.Id = id

	_, err = GetOne(w, r, id)
	if err != nil {
		return space, err
	}

	_, err = DB.Exec("UPDATE ad_space_auction.space SET name = ?, base_price = ? WHERE id = ?", space.Name, space.BasePrice, space.Id)
	if err != nil {
		return space, err
	}

	return space, err
}

func Delete(w http.ResponseWriter, r *http.Request, id string) (objects.SpaceEntity, error) {
	space, err := GetOne(w, r, id)
	if err != nil {
		return space, err
	}

	_, err = DB.Exec("DELETE FROM ad_space_auction.space WHERE id = ?", space.Id)
	if err != nil {
		return space, err
	}

	return space, err
}
