#!/bin/bash

## Download all tickers (Symbol, Name, MArker, PrimaryExchange)
instruments >tickers.tsv

## Load into new db
duckdb tickers.duckdb -c "DROP TABLE tickers; CREATE TABLE tickers AS select * from read_csv('tickers.tsv', auto_detect=true, header=true);"