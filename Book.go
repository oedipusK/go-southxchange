package southxchange

type BookOrder struct {
	Index  int     `json:"Index"`
	Amount float64 `json:"Amount"`
	Price  float64 `json:"Price"`
}

type Book struct {
	BuyOrders  []BookOrder `json:"BuyOrders"`
	SellOrders []BookOrder `json:"SellOrders"`
}