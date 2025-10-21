package main

import (
	"context"
	"fmt"
	"time"

	"github.com/SubhamShaww/realtime-bidding-system-engine/pkg/engine"
	"github.com/SubhamShaww/realtime-bidding-system-engine/pkg/utils"
)

func main() {
	fmt.Println("Realtime Bidding System Main Application")
	utils.InitPrometheus()
	bids := utils.GenerateBids(10)
	maxConcurrentBidders := 5
	metricsSize := 20

	// create a cancellable context for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// use mock bidder URL by default
	bidderURL := "http://localhost:8081/bid"
	engine.RunEngineWithBids(
		ctx,
		bids,
		maxConcurrentBidders,
		metricsSize,
		bidderURL,
	)
}
