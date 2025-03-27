// Stocks - Last Trade
// https://polygon.io/docs/stocks/get_v2_last_trade__stocksticker
// https://github.com/polygon-io/client-go/blob/master/rest/trades.go
package main

import (
	"context"
	"log"
	"os"
	"time"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)

func main() {

	// init client
	c := polygon.New(os.Getenv("POLYGON_API_KEY"))

	from, err := time.Parse("2006-01-02", "2025-03-20")
	if err != nil {
		log.Fatalf("Error parsing 'from' date: %v", err)
	}
	to, err := time.Parse("2006-01-02", "2025-03-26")
	if err != nil {
		log.Fatalf("Error parsing 'to' date: %v", err)
	}

	params := models.ListAggsParams{
		Ticker:     "AAPL",
		Multiplier: 1,
		Timespan:   "day",
		From:       models.Millis(from),
		To:         models.Millis(to),
	}.
		WithAdjusted(true).
		WithOrder(models.Order("asc")).
		WithLimit(120)

	iter := c.ListAggs(context.Background(), params)

	for iter.Next() {
		log.Print(iter.Item())
	}
	if iter.Err() != nil {
		log.Fatal(iter.Err())
	}

}
