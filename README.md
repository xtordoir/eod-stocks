# toolkit for eod stock data

## Get ticker list

`cmd/instrument/main.go`

```
go run cmd/instruments/main.go  >tickers.tsv
```

Format is tab separated, with the following colums:

```
Ticker  Name    Market  PrimaryExchange
```

## Load ticker list into duckDB

```
duckdb tickers.duckdb

CREATE TABLE tickers AS 
    select * from read_csv('tickers.tsv', auto_detect=true, header=true);
```

## Search for tickers and names

```
go run cmd/ticker/page/main.go
```

Creates a server on `localhost:8080` to explore the `tickers.duckdb` once it is created and loaded.