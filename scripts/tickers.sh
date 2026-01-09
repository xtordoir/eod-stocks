#!/bin/bash

## Download all tickers (Symbol, Name, MArker, PrimaryExchange)
instruments >tickers.tsv

## Load into new db
duckdb tickers.duckdb -c "DROP TABLE tickers; CREATE TABLE tickers AS select * from read_csv('tickers.tsv', auto_detect=true, header=true);"
## Set ticker status as not uptodate
duckdb tickers.duckdb -c "ALTER TABLE tickers ADD COLUMN uptodate boolean ; update tickers set uptodate = false;"

## Now download and create up to date commands
rm uptodate.sql
duckdb -noheader -csv tickers.duckdb "SELECT symbol FROM tickers WHERE uptodate = false" | \
  go run cmd/eod/main.go | \
  sed -ur "s/(.*)/UPDATE tickers SET uptodate = true WHERE symbol = '\1';/g" >>uptodate.sql

## And update status
cat uptodate.sql | duckdb tickers.duckdb 