package main

import (
	"fmt"

	"github.com/SubhamShaww/realtime-bidding-system-engine/pkg/engine"
	"github.com/SubhamShaww/realtime-bidding-system-engine/pkg/utils"
)

func main() {
	fmt.Println("Realtime Bidding System Main Application")
	bids := utils.GenerateBids(10)
	maxConcurrentBidders := 5
	metricsSize := 20
	engine.RunEngineWithBids(bids, maxConcurrentBidders, metricsSize)
}
