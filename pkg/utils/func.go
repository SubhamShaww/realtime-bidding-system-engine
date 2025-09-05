package utils

import "math/rand"

func GenerateBids(n int) []Bid {
	bids := make([]Bid, n)
	for i := range n {
		bids[i] = Bid{
			ID:       i + 1,
			Priority: rand.Intn(100), // Random priority between 0-99
		}
	}

	return bids
}
