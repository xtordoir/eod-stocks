package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"log"

	_ "github.com/duckdb/duckdb-go/v2"
)

var db *sql.DB

type ticker struct {
	ticker string
	name   string
}

func main() {
	fmt.Println(time.Now())

	var err error
	db, err = sql.Open("duckdb", "tickers.duckdb?access_mode=READ_ONLY")
	check(err)
	defer func() { check(db.Close()) }()

	fmt.Println(time.Now())
	queryAndPrint()
	fmt.Println(time.Now())
	queryAndPrint()
	fmt.Println(time.Now())

}

func queryAndPrint() {
	rows, err := db.QueryContext(
		context.Background(), `
		SELECT ticker, name
		FROM tickers
		WHERE ticker LIKE ?`,
		"Z%",
	)
	check(err)
	defer func() { check(rows.Close()) }()

	for rows.Next() {
		t := new(ticker)
		err := rows.Scan(&t.ticker, &t.name)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf(
			"Ticker: %s is %s\n",
			t.ticker, t.name,
		)
	}
	check(rows.Err())
}

func check(args ...any) {
	err := args[len(args)-1]
	if err != nil {
		panic(err)
	}
}
