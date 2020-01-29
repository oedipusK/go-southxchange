// Package SouthXchange is an implementation of the SouthXchange API in Golang.
package southxchange

import (
	"encoding/json"
	"errors"

	//"fmt"
	"net/http"
	"strconv"

	//"strings"
	"time"
	//"strings"
	"fmt"
)

const (
	API_BASE    = "https://www.southxchange.com/api" // SouthXchange API endpoint
	API_BASE_V2 = "https://www.southxchange.com/api/v2"
)

// New returns an instantiated SouthXchange struct
func New(apiKey, apiSecret, userAgent string) *SouthXchange {
	client := NewClient(apiKey, apiSecret, userAgent)
	return &SouthXchange{client}
}

// NewWithCustomHttpClient returns an instantiated SouthXchange struct with custom http client
func NewWithCustomHttpClient(apiKey, apiSecret, userAgent string, httpClient *http.Client) *SouthXchange {
	client := NewClientWithCustomHttpConfig(apiKey, apiSecret, userAgent, httpClient)
	return &SouthXchange{client}
}

// NewWithCustomTimeout returns an instantiated SouthXchange struct with custom timeout
func NewWithCustomTimeout(apiKey, apiSecret, userAgent string, timeout time.Duration) *SouthXchange {
	client := NewClientWithCustomTimeout(apiKey, apiSecret, userAgent, timeout)
	return &SouthXchange{client}
}

// handleErr gets JSON response from Cryptopia API en deal with error
func handleErr(r interface{}) error {
	switch v := r.(type) {
	case map[string]interface{}:
		error := r.(map[string]interface{})["error"]
		if error != nil {
			return errors.New(error.(string))
		}
		error_code := r.(map[string]interface{})["error_code"]
		if error_code != nil {
			return errors.New("Error: " + strconv.Itoa(error_code.(int)))
		}
	case []interface{}:
		return nil
	default:
		return fmt.Errorf("I don't know about type %T!\n", v)
	}

	return nil
}

// SouthXchange represent a SouthXchange client
type SouthXchange struct {
	client *client
}

// set enable/disable http request/response dump
func (o *SouthXchange) SetDebug(enable bool) {
	o.client.debug = enable
}

// GetMarketSummaries is used to get the last 24 hour summary of all active exchanges
func (b *SouthXchange) GetMarketSummaries() (marketSummaries []MarketSummary, err error) {
	response, err := b.client.do(API_BASE, "GET", "markets", nil, false)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &marketSummaries)
	return
}

// GetOpenOrders returns orders that you currently have opened.
func (b *SouthXchange) GetOpenOrders() (openOrders []Order, err error) {
	r, err := b.client.do(API_BASE, "POST", "listOrders", nil, true)
	if err != nil {
		return
	}
	err = json.Unmarshal(r, &openOrders)
	return
}

// GetOrder returns an order based on the orderCode
func (b *SouthXchange) GetOrder(orderCode string) (order Order, err error) {
	r, err := b.client.do(API_BASE_V2, "POST", "getOrder", map[string]string{"Code": orderCode}, true)
	if err != nil {
		return
	}

	err = json.Unmarshal(r, &order)
	return
}

// Account

// GetBalances is used to retrieve all balances from your account
func (o *SouthXchange) GetBalances() (balances []Balance, err error) {
	r, err := o.client.do(API_BASE, "POST", "listBalances", nil, true)
	if err != nil {
		return
	}
	err = json.Unmarshal(r, &balances)
	return
}

// GetDepositAddress is sed to generate or retrieve an address for a specific currency.
// currency a string literal for the currency (ie. BTC)
func (b *SouthXchange) GetDepositAddress(currency string) (address string, err error) {
	r, err := b.client.do(API_BASE, "POST", "generatenewaddress", map[string]string{"currency": currency}, true)
	if err != nil {
		return
	}
	address = string(r)
	return
}

// Withdraw is used to withdraw funds from your account.
// address string the address where to send the funds.
// currency string literal for the currency (ie. BTC)
// quantity float the quantity of coins to withdraw
// fee float the quantity of coins to withdraw
func (o *SouthXchange) Withdraw(address string, currency string, quantity float64) (withdraw WithdrawalInfo, err error) {
	r, err := o.client.do(API_BASE, "POST", "withdraw", map[string]string{
		"currency": currency,
		"address":  address,
		"amount":   strconv.FormatFloat(quantity, 'f', -1, 64),
	}, true)
	if err != nil {
		return
	}
	err = json.Unmarshal(r, &withdraw)
	return withdraw, err
}

// listing string
// reference string
// orderSide string enum
// amount float the quantity of coins to sell
// limitPrice float optional price in reference currency. ff nil then order is executed at market price
func (o *SouthXchange) PlaceOrder(listing string, reference string, orderSide OrderType, amount float64, limitPrice float64, marketPrice bool) (orderCode string, err error) {
	var params = make(map[string]string)
	params["listingCurrency"] = listing
	params["referenceCurrency"] = reference
	params["type"] = string(orderSide)
	params["amount"] = strconv.FormatFloat(amount, 'f', -1, 64)
	if !marketPrice {
		params["limitPrice"] = strconv.FormatFloat(limitPrice, 'f', -1, 64)
	}
	r, err := o.client.do(API_BASE, "POST", "placeOrder", params, true)
	if err != nil {
		return
	}
	orderCode = string(r)
	return
}

// gets and array containing all orders
func (o *SouthXchange) ListOrders() (orders []string, err error) {
	r, err := o.client.do(API_BASE, "POST", "listOrders", nil, true)
	if err != nil {
		return
	}
	for _, b := range r {
		orders = append(orders, string(b))
	}
	return
}

// GetTransactions is used to retrieve your transaction history
func (b *SouthXchange) GetTransactions(transactionType string, start uint64, limit uint32, sort string, desc bool) (transactions []Transaction, err error) {
	payload := make(map[string]string)
	if transactionType == "" {
		transactionType = "transactions"
	}
	if start > 0 {
		payload["PageIndex"] = strconv.FormatUint(uint64(start), 10)
	}
	if limit > 1000 {
		limit = 1000
	}
	if limit > 0 {
		payload["PageSize"] = strconv.FormatUint(uint64(limit), 10)
	}
	if sort == "" {
		sort = "Date"
	}
	payload["SortField"] = sort
	payload["Descending"] = strconv.FormatBool(desc)
	payload["TransactionType"] = transactionType
	r, err := b.client.do(API_BASE, "POST", "listTransactions", payload, true)
	if err != nil {
		return
	}
	var res struct {
		TotalElements int
		Result        []Transaction
	}
	err = json.Unmarshal(r, &res)
	return res.Result, err
}
