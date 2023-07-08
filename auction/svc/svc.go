package svc

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"auction/objects"

	"github.com/google/uuid"
)

var DB *sql.DB

func GetAll(w http.ResponseWriter, r *http.Request) (objects.GetAllAuctionResp, error) {
	var resp = objects.GetAllAuctionResp{}

	rows, err := DB.Query("SELECT * FROM ad_space_auction.auction")
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	var auctions = []objects.AuctionEntity{}
	for rows.Next() {
		var auction objects.AuctionEntity
		err := rows.Scan(&auction.Id, &auction.SpaceId, &auction.StartsOn, &auction.EndsOn)
		if err != nil {
			return resp, err
		}
		auctions = append(auctions, auction)
	}

	if err := rows.Err(); err != nil {
		return resp, err
	}

	resp.Data = auctions
	return resp, err
}

func GetOne(w http.ResponseWriter, r *http.Request, id string) (objects.GetOneAuctionResp, error) {
	row := DB.QueryRow("SELECT * FROM ad_space_auction.auction WHERE id = ?", id)

	var resp objects.GetOneAuctionResp
	var auction objects.AuctionEntity

	err := row.Scan(&auction.Id, &auction.SpaceId, &auction.StartsOn, &auction.EndsOn)
	if err == sql.ErrNoRows {
		http.Error(w, "auction entity not found", http.StatusNotFound)
		return resp, err
	}
	if err != nil {
		return resp, err
	}

	space, err := GetSpaceEntity(auction.SpaceId)
	if err != nil {
		log.Println(err)
		err = errors.New("space_id doesn't exist for given auction")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return resp, err
	}

	existingBids, err := GetAllAuctionBids(auction.Id)
	if err != nil {
		log.Println(err)
		err = errors.New("can't get existing auction bids")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return resp, err
	}

	resp.AuctionEntity = auction
	resp.SpaceName = space.Name
	resp.SpaceBasePrice = space.BasePrice
	resp.AuctionBids = existingBids

	return resp, err
}

func Create(w http.ResponseWriter, r *http.Request) (objects.AuctionEntity, error) {
	var reqBody objects.CreateAuctionReq
	var auction = objects.AuctionEntity{}
	id := uuid.NewString()
	startsOn := time.Now()

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return auction, err
	}

	if reqBody.SpaceId == "" {
		err = errors.New("space_id shouldn't be empty")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return auction, err
	} else if reqBody.EndsOn == "" {
		err = errors.New("ends_on shouldn't be empty")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return auction, err
	}

	endsOn, err := time.Parse(time.RFC3339, reqBody.EndsOn)
	if err != nil {
		err = errors.New("ends_on should be in RFC3339 time format")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return auction, err
	}

	if endsOn.Compare(startsOn) < 1 {
		err = errors.New("ends_on should be more than current time")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return auction, err
	}

	space, err := GetSpaceEntity(reqBody.SpaceId)
	if err != nil {
		log.Println(err)
		err = errors.New("space_id doesn't exist")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return auction, err
	}

	auction.Id = id
	auction.SpaceId = space.Id
	auction.StartsOn = startsOn
	auction.EndsOn = endsOn

	_, err = DB.Exec("INSERT INTO ad_space_auction.auction (id, space_id, ends_on) VALUES (?, ?, ?)", auction.Id, auction.SpaceId, auction.EndsOn)
	if err != nil {
		return auction, err
	}

	return auction, err
}

func GetSpaceEntity(spaceId string) (objects.SpaceEntity, error) {
	supplySvcHost := os.Getenv("SUPPLY_SVC_HOST")
	supplySvcPort := os.Getenv("SUPPLY_SVC_PORT")

	var space objects.SpaceEntity

	resp, err := http.Get("http://" + supplySvcHost + ":" + supplySvcPort + "/space/" + spaceId)
	if err != nil {
		return space, err
	}

	if statusCodeBit := resp.StatusCode / 100; statusCodeBit != 2 && statusCodeBit != 3 {
		err = fmt.Errorf("api failed with status code %d", resp.StatusCode)
		return space, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&space)
	if err != nil {
		return space, err
	}

	return space, err
}

func GetBidderEntity(bidderId string) (objects.BidderEntity, error) {
	demandSvcHost := os.Getenv("DEMAND_SVC_HOST")
	demandSvcPort := os.Getenv("DEMAND_SVC_PORT")

	var bidder objects.BidderEntity

	resp, err := http.Get("http://" + demandSvcHost + ":" + demandSvcPort + "/bidder/" + bidderId)
	if err != nil {
		return bidder, err
	}

	if statusCodeBit := resp.StatusCode / 100; statusCodeBit != 2 && statusCodeBit != 3 {
		err = fmt.Errorf("api failed with status code %d", resp.StatusCode)
		return bidder, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&bidder)
	if err != nil {
		return bidder, err
	}

	return bidder, err
}

func GetAllAuctionBids(auctionId string) ([]objects.AuctionBidEntity, error) {
	var resp = []objects.AuctionBidEntity{}

	rows, err := DB.Query("SELECT * FROM ad_space_auction.auction_bid WHERE auction_id = ? ORDER BY bid_price DESC", auctionId)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		var auctionBid objects.AuctionBidEntity
		err := rows.Scan(&auctionBid.Id, &auctionBid.AuctionId, &auctionBid.BidderId, &auctionBid.BidPrice, &auctionBid.BidOn)
		if err != nil {
			return resp, err
		}
		resp = append(resp, auctionBid)
	}

	if err := rows.Err(); err != nil {
		return resp, err
	}

	return resp, err
}

func CreateAuctionBid(w http.ResponseWriter, r *http.Request, auction_id string) (objects.AuctionBidEntity, error) {
	var reqBody objects.CreateAuctionBidReq
	var auctionBid = objects.AuctionBidEntity{}
	id := uuid.NewString()
	bid_on := time.Now()

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return auctionBid, err
	}

	if reqBody.BidderId == "" {
		err = errors.New("bidder_id shouldn't be empty")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return auctionBid, err
	} else if reqBody.BidPrice <= 0 {
		err = errors.New("bid_price should be more than zero")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return auctionBid, err
	}

	bidder, err := GetBidderEntity(reqBody.BidderId)
	if err != nil {
		log.Println(err)
		err = errors.New("bidder_id doesn't exist")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return auctionBid, err
	}

	auction, err := GetOne(w, r, auction_id)
	if err != nil {
		return auctionBid, err
	}

	if bid_on.Compare(auction.EndsOn) > -1 {
		err = errors.New("auction is closed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return auctionBid, err
	}

	if reqBody.BidPrice <= auction.SpaceBasePrice {
		err = errors.New("bid_price should be more than space_base_price")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return auctionBid, err
	}

	if len(auction.AuctionBids) > 0 && reqBody.BidPrice <= auction.AuctionBids[0].BidPrice {
		err = fmt.Errorf("bid_price should be more than last bid price %v", auction.AuctionBids[0].BidPrice)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return auctionBid, err
	}

	auctionBid.Id = id
	auctionBid.AuctionId = auction.Id
	auctionBid.BidderId = bidder.Id
	auctionBid.BidPrice = reqBody.BidPrice
	auctionBid.BidOn = bid_on

	_, err = DB.Exec("INSERT INTO ad_space_auction.auction_bid (id, auction_id, bidder_id, bid_price, bid_on) VALUES (?, ?, ?, ?, ?)", auctionBid.Id, auctionBid.AuctionId, auctionBid.BidderId, auctionBid.BidPrice, auctionBid.BidOn)
	if err != nil {
		return auctionBid, err
	}

	return auctionBid, err
}
