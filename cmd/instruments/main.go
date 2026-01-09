package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	massive "github.com/massive-com/client-go/v2/rest"
	"github.com/massive-com/client-go/v2/rest/models"
)

func main() {

	c := massive.New(os.Getenv("POLYGON_API_KEY"))

	params := models.ListTickersParams{}.
		WithMarket(models.AssetClass("stocks")).
		WithActive(true).
		WithOrder(models.Order("asc")).
		WithLimit(120).
		WithSort(models.Sort("ticker"))

	iter := c.ListTickers(context.Background(), params)

	fmt.Println("Symbol" + "\t" + "Name" + "\t" + "Market" + "\t" + "PrimaryExchange")

	for iter.Next() {
		ticker := iter.Item()
		fmt.Println(ticker.Ticker + "\t" + ticker.Name + "\t" + ticker.Market + "\t" + ticker.PrimaryExchange)
		// Calling Sleep method
		time.Sleep(100 * time.Millisecond)
	}
	if iter.Err() != nil {
		log.Fatal(iter.Err())
	}
}
