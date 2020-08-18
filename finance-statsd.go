package main

import (
	"fmt"

	"github.com/piquette/finance-go/quote"
)

func main() {
	symbols := []string{"AAPL", "TSLA", "MSFT"}
	iter := quote.List(symbols)

	// Iterate over results. Will exit upon any error.
	for iter.Next() {
		q := iter.Quote()
		fmt.Println(q.Symbol, ": ", q.Bid)
	}

	// Catch an error, if there was one.
	if iter.Err() != nil {
		// Uh-oh!
		panic(iter.Err())
	}
}
