package southxchange

type MarketPrice struct {
	Bid           string `json:"Bid"`
	Ask           string `json:"Ask"`
	Last          string `json:"Last"`
	Variation24Hr string `json:"Variation24Hr"`
	Volume24Hr    string `json:"Volume24Hr"`
}
