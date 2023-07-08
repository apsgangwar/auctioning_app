package objects

import "time"

type AuctionEntity struct {
	Id       string    `json:"id"`
	SpaceId  string    `json:"space_id"`
	StartsOn time.Time `json:"starts_on"`
	EndsOn   time.Time `json:"ends_on"`
}

type CreateAuctionReq struct {
	SpaceId string `json:"space_id"`
	EndsOn  string `json:"ends_on"`
}

type GetAllAuctionResp struct {
	Data []AuctionEntity `json:"data"`
}

type GetOneAuctionResp struct {
	AuctionEntity
	SpaceName      string             `json:"space_name"`
	SpaceBasePrice float64            `json:"space_base_price"`
	AuctionBids    []AuctionBidEntity `json:"auction_bids"`
}

type SpaceEntity struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	BasePrice float64 `json:"base_price"`
}

type BidderEntity struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AuctionBidEntity struct {
	Id        string    `json:"id"`
	AuctionId string    `json:"auction_id"`
	BidderId  string    `json:"bidder_id"`
	BidPrice  float64   `json:"bid_price"`
	BidOn     time.Time `json:"bid_on"`
}

type CreateAuctionBidReq struct {
	BidderId string  `json:"bidder_id"`
	BidPrice float64 `json:"bid_price"`
}
