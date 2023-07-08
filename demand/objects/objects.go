package objects

type BidderEntity struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	// Email     string `json:"email"`
	// ContactNo string `json:"contact_no"`
	// Address   string `json:"address"`
}

type GetAllBidderResp struct {
	Data []BidderEntity `json:"data"`
}
