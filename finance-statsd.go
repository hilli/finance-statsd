package main

import (
	"fmt"

	"github.com/piquette/finance-go/equity"
)

// Uses https://piquette.io/projects/finance-go/ https://github.com/piquette/finance-go

func main() {
	symbols := []string{"AAPL", "TSLA", "MSFT"}
	// iter := quote.List(symbols)
	iter := equity.List(symbols)

	// Iterate over results. Will exit upon any error.
	for iter.Next() {
		// q := iter.Quote()
		q := iter.Equity()
		fmt.Println(q.Symbol, "(", q.ShortName, "): Bid: ", q.Bid, " Ask: ", q.Ask, "Price:", q.RegularMarketPrice)
		fmt.Println("High:", q.RegularMarketDayHigh, "Low:", q.RegularMarketDayLow, "Close:", q.RegularMarketPreviousClose)
		fmt.Println()

	}

	// Catch an error, if there was one.
	if iter.Err() != nil {
		// Uh-oh!
		panic(iter.Err())
	}
}
