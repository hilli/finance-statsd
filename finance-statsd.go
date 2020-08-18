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

// Uses https://piquette.io/projects/finance-go/ https://github.com/piquette/finance-go

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
		ci = 10 // ci probably unset - But if set to 0 then don't DoS the bridge
	}
	collectionInterval := time.Duration(ci) * time.Second

	symbols := strings.Split(os.Getenv("SYMBOLS"), ",")

	for {
		iter := equity.List(symbols)

		// Iterate over results. Will exit upon any error.
		for iter.Next() {
			q := iter.Equity()
			if os.Getenv("DEBUG") != "" {
				fmt.Println(q.Symbol, "(", q.ShortName, "): Bid: ", q.Bid, " Ask: ", q.Ask, "Price:", q.RegularMarketPrice)
				fmt.Printf("     High: %.2f Low: %.2f Close: %.2f Post: %.2f\n", q.RegularMarketDayHigh, q.RegularMarketDayLow, q.RegularMarketPreviousClose, q.RegularMarketPrice+q.PreMarketChange)
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
