// Stocks - Last Trade
// https://polygon.io/docs/stocks/get_v2_last_trade__stocksticker
// https://github.com/polygon-io/client-go/blob/master/rest/trades.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)

func main() {

	// init client
	c := polygon.New(os.Getenv("POLYGON_API_KEY"))

	// read the symbols line by line from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		symbol := scanner.Text()
		downloadEOD(symbol, *c)

		time.Sleep(15 * time.Second)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func downloadEOD(symbol string, c polygon.Client) {

	from, err := time.Parse("2006-01-02", "2025-01-08")
	if err != nil {
		log.Fatalf("Error parsing 'from' date: %v", err)
	}
	to, err := time.Parse("2006-01-02", "2026-01-09")
	if err != nil {
		log.Fatalf("Error parsing 'to' date: %v", err)
	}

	params := models.ListAggsParams{
		Ticker:     symbol,
		Multiplier: 1,
		Timespan:   "day",
		From:       models.Millis(from),
		To:         models.Millis(to),
	}.
		WithAdjusted(true).
		WithOrder(models.Order("asc")).
		WithLimit(600)

	iter := c.ListAggs(context.Background(), params)
	f, err := os.Create("data/" + symbol + ".tsv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		f.Close()
		_, err := fmt.Println(symbol)
		if err != nil {
			log.Println(err)
		}
		log.Println(symbol)
	}()

	fmt.Fprintf(f, "Date\tOpen\tHigh\tLow\tClose\tVolume\n")

	for iter.Next() {
		agg := iter.Item()
		fmt.Fprintf(f, "%s\t%.2f\t%.2f\t%.2f\t%.2f\t%.0f\n", time.Time(agg.Timestamp).Format("2006-01-02"), agg.Open, agg.High, agg.Low, agg.Close, agg.Volume)
	}
}

/***
	os := gohff.Overshoot{
		Instrument:    "AAPL",
		Scale:         5.0,
		Direction:     1,
		StartExtremum: 1.0,
		PeakExtremum:  1.0,
		Current:       1.0,
	}

	n1 := 0
	n2 := 0
	n3 := 0
	n4 := 0
	n5 := 0
	p1 := 0
	p2 := 0
	p3 := 0
	p4 := 0
	p5 := 0
	sumP := 0.0
	nP := 0.0
	sumN := 0.0
	nN := 0.0
	iter.Next()
	agg := iter.Item()
	os.Reset(agg.Open)
	os = os.Update(agg.Close)

	for iter.Next() {
		agg := iter.Item()
		//log.Println(time.Time(agg.Timestamp))

		//log.Print(iter.Item())
		n := os.Update(agg.Open)
		if n.MaxOS()*os.MaxOS() < 0 {
			log.Printf("%+v %.2f\n", os.MaxOS(), os.PeakExtremum)
			if os.MaxOS() < 0 {
				sumP = sumP + os.MaxOS()
				nP = nP + 1
				if os.MaxOS() > -2 {
					p1 = p1 + 1
				} else if os.MaxOS() > -3 {
					p2 = p2 + 1
				} else if os.MaxOS() > -4 {
					p3 = p3 + 1
				} else if os.MaxOS() > -5 {
					p4 = p4 + 1
				} else {
					p5 = p5 + 1
				}
			}
			if os.MaxOS() > 0 {
				sumN = sumN + os.MaxOS()
				nN = nN + 1
				if os.MaxOS() < 2 {
					n1 = n1 + 1
				} else if os.MaxOS() < 3 {
					n2 = n2 + 1
				} else if os.MaxOS() < 4 {
					n3 = n3 + 1
				} else if os.MaxOS() < 5 {
					n4 = n4 + 1
				} else {
					n5 = n5 + 1
				}
			}
		}
		w := n.Update(agg.Close)
		if n.MaxOS()*w.MaxOS() < 0 {
			log.Printf("%+v %.2f\n", n.MaxOS(), n.PeakExtremum)
			if n.MaxOS() < 0 {
				sumP = sumP + n.MaxOS()
				nP = nP + 1
				if n.MaxOS() > -2 {
					p1 = p1 + 1
				} else if n.MaxOS() > -3 {
					p2 = p2 + 1
				} else if n.MaxOS() > -4 {
					p3 = p3 + 1
				} else if n.MaxOS() > -5 {
					p4 = p4 + 1
				} else {
					p5 = p5 + 1
				}
			}
			if n.MaxOS() > 0 {
				sumN = sumN + n.MaxOS()
				nN = nN + 1
				if n.MaxOS() < 2 {
					n1 = n1 + 1
				} else if n.MaxOS() < 3 {
					n2 = n2 + 1
				} else if n.MaxOS() < 4 {
					n3 = n3 + 1
				} else if n.MaxOS() < 5 {
					n4 = n4 + 1
				} else {
					n5 = n5 + 1
				}
			}
		}
		os = w
		//log.Printf("%+v\n", os.MaxOS())
	}
	if iter.Err() != nil {
		log.Fatal(iter.Err())
	}

	log.Printf("%d %d %d %d %d : %.2f\n", n1, n2, n3, n4, n5, sumN/nN)
	log.Printf("%d %d %d %d %d : %.2f\n", p1, p2, p3, p4, p5, sumP/nP)
}
***/
