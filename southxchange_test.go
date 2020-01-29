package southxchange

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}
}

func TestOrderPlacement(t *testing.T) {
	southxchange := New(os.Getenv("SOUTH_API_KEY"), os.Getenv("SOUTH_API_SECRET"), "user-agent 1.0")
	res, err := southxchange.PlaceOrder("polis", "btc", "sell", 5, 0.0, true)
	assert.Nil(t, err)
	fmt.Println(res)
}

func TestListOrders(t *testing.T) {
	southxchange := New(os.Getenv("SOUTH_API_KEY"), os.Getenv("SOUTH_API_SECRET"), "user-agent 1.0")
	res, err := southxchange.ListOrders()
	assert.Nil(t, err)
	fmt.Println(res)
}
