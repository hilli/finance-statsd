package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	_ "github.com/joho/godotenv/autoload"
	"github.com/piquette/finance-go/equity"
)

// Uses https://piquette.io/projects/finance-go/ from https://github.com/piquette/finance-go - Thanks piquette

func main() {
	// Setup
	statsdEndpoint := os.Getenv("STATSD_ENDPOINT")
	if statsdEndpoint == "" {
		statsdEndpoint = "127.0.0.1:8125"
	}
	statsd, err := statsd.New(statsdEndpoint)
	statsd.Namespace = "stocks."
	if err != nil {
		log.Fatal(err)
	}

	// Setup collection interval
	ci, _ := strconv.ParseUint(os.Getenv("COLLECTION_INTERVAL"), 10, 32)
	if ci == 0 {
		ci = 60
	}
	collectionInterval := time.Duration(ci) * time.Second

	symbols := strings.Split(os.Getenv("SYMBOLS"), ",")
	if len(symbols) == 0 || symbols[0] == "" {
		panic("Please set the environment variable SYMBOLS to your preferred stock ticker names, comma seperated. Ie SYMBOLS=AAPL,TSLA")
	}

	// Collection loop
	for {
		iter := equity.List(symbols)

		// Iterate over results. Will exit upon any error.
		for iter.Next() {
			q := iter.Equity()
			if os.Getenv("DEBUG") != "" {
				fmt.Printf("%s (%s): Bid: %.2f Ask: %.2f Price: %.2f High: %.2f Low: %.2f Close: %.2f Post: %.2f\n", q.Symbol, q.ShortName, q.Bid, q.Ask, q.RegularMarketPrice, q.RegularMarketDayHigh, q.RegularMarketDayLow, q.RegularMarketPreviousClose, q.RegularMarketPrice+q.PreMarketChange)
			}
			statsd.Gauge("bid", q.Bid, []string{"symbol:" + q.Symbol}, 1)
			statsd.Gauge("ask", q.Ask, []string{"symbol:" + q.Symbol}, 1)
			statsd.Gauge("price", q.RegularMarketPrice, []string{"symbol:" + q.Symbol}, 1)
			statsd.Gauge("high", q.RegularMarketDayHigh, []string{"symbol:" + q.Symbol}, 1)
			statsd.Gauge("low", q.RegularMarketDayLow, []string{"symbol:" + q.Symbol}, 1)
			statsd.Gauge("prev_close", q.RegularMarketPreviousClose, []string{"symbol:" + q.Symbol}, 1)
			statsd.Gauge("post", q.RegularMarketPrice+q.PreMarketChange, []string{"symbol:" + q.Symbol}, 1)
		}

		// Catch an error, if there was one.
		if iter.Err() != nil {
			// Uh-oh!
			panic(iter.Err())
		}
		time.Sleep(collectionInterval)
	}
}
