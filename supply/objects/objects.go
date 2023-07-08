package objects

type SpaceEntity struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	BasePrice float64 `json:"base_price"`
	// Website   string  `json:"website"`
	// Size      string  `json:"size"`
	// Position  string  `json:"position"`
}

type GetAllSpaceResp struct {
	Data []SpaceEntity `json:"data"`
}
