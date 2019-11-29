package southxchange

type Order struct {
	Code              string
	Type              string
	Amount            float64
	OriginalAmount    float64
	LimitPrice        float64
	ListingCurrency   string
	ReferenceCurrency string
}

type PlacedOrderResponse struct {
	OrderCode string
}

type OrderType string

const (
	Buy  OrderType = "buy"
	Sell OrderType = "sell"
)
