package engine

import (
	"math/rand"
)

func generateBids(n int) []Bid {
	bids := make([]Bid, n)
	for i := 0; i < n; i++ {
		bids[i] = Bid{
			ID:       i + 1,
			Priority: rand.Intn(100), // Random priority between 0-99
		}
	}

	return bids
}
