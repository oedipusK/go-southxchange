package southxchange

type MarketPrice struct {
	Bid           float64 `json:"Bid"`
	Ask           float64 `json:"Ask"`
	Last          float64 `json:"Last"`
	Variation24Hr float64 `json:"Variation24Hr"`
	Volume24Hr    float64 `json:"Volume24Hr"`
}
